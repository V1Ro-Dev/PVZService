package repository_test

import (
	"context"
	"database/sql"
	"errors"

	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"pvz/internal/models"
	"pvz/internal/repository"
	"pvz/internal/repository/mocks"
)

func TestCreateReception(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresReceptionRepository{Db: db}

	reception := models.Reception{
		Id:       uuid.New(),
		DateTime: time.Now().Truncate(time.Millisecond),
		PvzId:    uuid.New(),
		Status:   models.InProgress,
	}

	tests := []struct {
		name         string
		setupMock    func()
		expectedErr  bool
		expectedText string
	}{
		{
			name: "successfully creates reception",
			setupMock: func() {
				mock.ExpectExec("insert into reception").
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
		},
		{
			name: "no rows affected",
			setupMock: func() {
				mock.ExpectExec("insert into reception").
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr:  true,
			expectedText: "one reception was not closed or non-existing pvzId",
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectExec("insert into reception").
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnError(errors.New("db error"))
			},
			expectedErr:  true,
			expectedText: "unable to create reception",
		},
		{
			name: "pg error",
			setupMock: func() {
				pgErr := &pgconn.PgError{
					Message: "violates foreign key constraint",
					Detail:  "Key (pvz_id)=(...) is not present in table pvz.",
					Where:   "SQL insert",
				}
				mock.ExpectExec("insert into reception").
					WithArgs(reception.Id, reception.DateTime, reception.PvzId, reception.Status).
					WillReturnError(pgErr)
			},
			expectedErr:  true,
			expectedText: "SQL Error: violates foreign key constraint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.CreateReception(context.Background(), reception)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectedErr && tt.expectedText != "" && err != nil && !strings.Contains(err.Error(), tt.expectedText) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedText, err.Error())
			}
		})
	}
}

func TestGetOpenReception(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresReceptionRepository{Db: db}

	pvzId := uuid.New()
	reception := models.Reception{
		Id:       uuid.New(),
		DateTime: time.Now().Truncate(time.Millisecond),
		PvzId:    pvzId,
		Status:   models.InProgress,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
		expected    models.Reception
	}{
		{
			name: "successfully retrieves open reception",
			setupMock: func() {
				mock.ExpectQuery("select id, reception_datetime, pvz_id, status").
					WithArgs(pvzId, models.InProgress).
					WillReturnRows(sqlmock.NewRows([]string{"id", "reception_datetime", "pvz_id", "status"}).
						AddRow(reception.Id, reception.DateTime, reception.PvzId, reception.Status))
			},
			expectedErr: false,
			expected:    reception,
		},
		{
			name: "no open reception found",
			setupMock: func() {
				mock.ExpectQuery("select id, reception_datetime, pvz_id, status").
					WithArgs(pvzId, models.InProgress).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: false,
			expected:    models.Reception{},
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("select id, reception_datetime, pvz_id, status").
					WithArgs(pvzId, models.InProgress).
					WillReturnError(errors.New("db error"))
			},
			expectedErr: true,
			expected:    models.Reception{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.GetOpenReception(context.Background(), pvzId)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAddProduct(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresReceptionRepository{Db: db}

	product := models.Product{
		Id:          uuid.New(),
		DateTime:    time.Now(),
		ProductType: "TypeA",
		ReceptionId: uuid.New(),
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
	}{
		{
			name: "successfully adds product to open reception",
			setupMock: func() {
				mock.ExpectExec("insert into product").
					WithArgs(product.Id, product.DateTime, product.ProductType, product.ReceptionId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
		},
		{
			name: "query error while adding product",
			setupMock: func() {
				mock.ExpectExec("insert into product").
					WithArgs(product.Id, product.DateTime, product.ProductType, product.ReceptionId).
					WillReturnError(errors.New("db error"))
			},
			expectedErr: true,
		},
		{
			name: "pg error",
			setupMock: func() {
				mock.ExpectExec("insert into product").
					WithArgs(product.Id, product.DateTime, product.ProductType, product.ReceptionId).
					WillReturnError(&pgconn.PgError{
						Message: "some weird SQL Error",
						Detail:  "Super Mega Detailed error",
						Where:   "In closing reception",
					})
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.AddProduct(context.Background(), product)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRemoveProduct(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresReceptionRepository{Db: db}

	receptionId := uuid.New()
	product := models.Product{
		Id:          uuid.New(),
		DateTime:    time.Now(),
		ProductType: "TypeA",
		ReceptionId: receptionId,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
	}{
		{
			name: "successfully removes last product from reception",
			setupMock: func() {
				mock.ExpectQuery("delete from product").
					WithArgs(receptionId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "received_at", "type", "reception_id"}).
						AddRow(product.Id, product.DateTime, product.ProductType, product.ReceptionId))
			},
			expectedErr: false,
		},
		{
			name: "no products to remove",
			setupMock: func() {
				mock.ExpectQuery("delete from product").
					WithArgs(receptionId).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: true,
		},
		{
			name: "query error while removing product",
			setupMock: func() {
				mock.ExpectQuery("delete from product").
					WithArgs(receptionId).
					WillReturnError(errors.New("db error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.RemoveProduct(context.Background(), receptionId)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCloseReception(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresReceptionRepository{Db: db}

	reception := models.Reception{
		Id:       uuid.New(),
		DateTime: time.Now(),
		PvzId:    uuid.New(),
		Status:   models.InProgress,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
	}{
		{
			name: "successfully closes reception",
			setupMock: func() {
				mock.ExpectExec("update reception set status").
					WithArgs(reception.Id, models.Closed).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
		},
		{
			name: "query error while closing reception",
			setupMock: func() {
				mock.ExpectExec("update reception set status").
					WithArgs(reception.Id, models.Closed).
					WillReturnError(errors.New("db error"))
			},
			expectedErr: true,
		},
		{
			name: "pg error",
			setupMock: func() {
				mock.ExpectExec("update reception set status").
					WithArgs(reception.Id, models.Closed).
					WillReturnError(&pgconn.PgError{
						Message: "some weird SQL Error",
						Detail:  "Super Mega Detailed error",
						Where:   "In closing reception",
					})
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.CloseReception(context.Background(), reception)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

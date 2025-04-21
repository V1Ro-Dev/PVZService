package repository_test

import (
	"context"
	"errors"
	"pvz/internal/delivery/forms"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"pvz/internal/models"
	"pvz/internal/repository"
	"pvz/internal/repository/mocks"
)

func TestCreatePvz(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresPvzRepository{Db: db}

	type args struct {
		ctx     context.Context
		pvzData models.Pvz
	}

	id := uuid.New()
	date := time.Now().Truncate(time.Millisecond)
	city := "Москва"

	tests := []struct {
		name      string
		args      args
		mockQuery func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				pvzData: models.Pvz{
					Id:               id,
					RegistrationDate: date,
					City:             city,
				},
			},
			mockQuery: func() {
				mock.ExpectExec(`insert into pvz`).
					WithArgs(id, date, city).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "sql error",
			args: args{
				ctx: context.Background(),
				pvzData: models.Pvz{
					Id:               id,
					RegistrationDate: date,
					City:             city,
				},
			},
			mockQuery: func() {
				mock.ExpectExec(`insert into pvz`).
					WithArgs(id, date, city).
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			err := repo.CreatePvz(tt.args.ctx, tt.args.pvzData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPvzInfo(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresPvzRepository{Db: db}

	start := time.Now().Truncate(time.Millisecond)
	end := time.Now().Truncate(time.Millisecond)

	form := forms.GetPvzInfoForm{
		StartDate: start,
		EndDate:   end,
		Page:      1,
		Limit:     10,
	}

	pvzId := uuid.New()
	receptionId := uuid.New()
	productId := uuid.New()

	tests := []struct {
		name      string
		mockQuery func()
		wantErr   bool
	}{
		{
			name: "success",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{
					"pvz_id", "pvz_registration_date", "pvz_city",
					"id", "reception_datetime", "status", "pvz_id",
					"product_id", "received_at", "type", "reception_id",
				}).AddRow(
					pvzId, start, "Москва",
					receptionId, start, string(models.InProgress), pvzId,
					productId, end, "обувь", receptionId,
				)
				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzInfoQuery)).
					WithArgs(form.StartDate, form.EndDate, form.Limit, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "query error",
			mockQuery: func() {
				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzInfoQuery)).
					WithArgs(form.StartDate, form.EndDate, form.Limit, 0).
					WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
		{
			name: "scan error",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{
					"pvz_id", "pvz_registration_date", "pvz_city",
					"id", "reception_datetime", "status", "pvz_id",
					"product_id", "received_at", "type", "reception_id",
				}).AddRow(
					"invalid-uuid", start, "Москва",
					receptionId, start, string(models.InProgress), pvzId,
					productId, end, "обувь", receptionId,
				)
				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzInfoQuery)).
					WithArgs(form.StartDate, form.EndDate, form.Limit, 0).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			_, err := repo.GetPvzInfo(context.Background(), form)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

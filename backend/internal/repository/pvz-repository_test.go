package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/forms"
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
			name: "ok",
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
		{
			name: "pg error",
			args: args{
				ctx: context.Background(),
				pvzData: models.Pvz{
					Id:               id,
					RegistrationDate: date,
					City:             city,
				},
			},
			mockQuery: func() {
				pgErr := &pgconn.PgError{
					Message: "duplicate key",
					Detail:  "Key (id)=(...) already exists.",
					Where:   "SQL statement",
				}
				mock.ExpectExec(`insert into pvz`).
					WithArgs(id, date, city).
					WillReturnError(pgErr)
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
	productType := "обувь"
	city := "Москва"

	tests := []struct {
		name      string
		mockQuery func()
		wantErr   bool
	}{

		{
			name: "ok",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{
					"id", "registration_date", "city",
					"id", "reception_datetime", "status", "pvz_id",
					"id", "received_at", "type", "reception_id",
				}).AddRow(
					pvzId, start, city,
					receptionId, start, string(models.InProgress), pvzId,
					productId, end, productType, receptionId,
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
					"id", "registration_date", "city",
					"id", "reception_datetime", "status", "pvz_id",
					"id", "received_at", "type", "reception_id",
				}).AddRow(
					"invalid-uuid", start, city,
					receptionId, start, string(models.InProgress), pvzId,
					productId, end, productType, receptionId,
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

func TestGetPvzList(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresPvzRepository{Db: db}

	idFirst := uuid.New()
	dateFirst := time.Now().Truncate(time.Millisecond)
	cityFirst := "Москва"
	idSecond := uuid.New()
	dateSecond := time.Now().Truncate(time.Millisecond)
	citySecond := "Казань"

	tests := []struct {
		name      string
		mockQuery func()
		wantErr   bool
		want      []models.Pvz
	}{
		{
			name: "ok",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow(idFirst.String(), dateFirst, cityFirst).
					AddRow(idSecond.String(), dateSecond, citySecond)

				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzListQuery)).
					WillReturnRows(rows)
			},
			wantErr: false,
			want: []models.Pvz{
				{Id: idFirst, RegistrationDate: dateFirst, City: cityFirst},
				{Id: idSecond, RegistrationDate: dateSecond, City: citySecond},
			},
		},
		{
			name: "query error",
			mockQuery: func() {
				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzListQuery)).
					WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "scan error",
			mockQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
					AddRow("invalid-uuid", dateFirst, cityFirst)

				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzListQuery)).
					WillReturnRows(rows)
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "empty result",
			mockQuery: func() {
				mock.ExpectQuery(regexp.QuoteMeta(repository.GetPvzListQuery)).WillReturnRows(sqlmock.NewRows([]string{}))
			},
			wantErr: false,
			want:    []models.Pvz(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()

			got, err := repo.GetPvzList(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

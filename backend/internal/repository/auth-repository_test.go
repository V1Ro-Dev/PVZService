package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"pvz/internal/models"
	"pvz/internal/repository"
	"pvz/internal/repository/mocks"
)

func TestCreateUser(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresUserRepository{Db: db}
	user := models.User{
		Id:       uuid.New().String(),
		Email:    "abobus@mail.ru",
		Password: "superMegaHashUnrealNoWayReally?HashedPassword",
		Salt:     "saltySalt",
		Role:     string(models.Client),
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`insert into "user" (id, email, password, salt, role) values ($1, $2, $3, $4, $5)`)).
					WithArgs(user.Id, user.Email, user.Password, user.Salt, user.Role).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
		},
		{
			name: "error on sql execution",
			setupMock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`insert into "user" (id, email, password, salt, role) values ($1, $2, $3, $4, $5)`)).
					WithArgs(user.Id, user.Email, user.Password, user.Salt, user.Role).
					WillReturnError(fmt.Errorf("SQL error"))
			},
			expectedErr: true,
		},
		{
			name: "error on duplicate email",
			setupMock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`insert into "user" (id, email, password, salt, role) values ($1, $2, $3, $4, $5)`)).
					WithArgs(user.Id, user.Email, user.Password, user.Salt, user.Role).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.CreateUser(context.Background(), user)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIsUserExist(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresUserRepository{Db: db}
	email := "abobus@mail.ru"
	userId := uuid.New()

	tests := []struct {
		name        string
		setupMock   func()
		expected    bool
		expectedErr bool
	}{
		{
			name: "user exists",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(userId)
				mock.ExpectQuery(regexp.QuoteMeta(`select id from "user" where email = $1`)).
					WithArgs(email).
					WillReturnRows(rows)
			},
			expected:    true,
			expectedErr: false,
		},
		{
			name: "user does not exist",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select id from "user" where email = $1`)).
					WithArgs(email).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    false,
			expectedErr: false,
		},
		{
			name: "unexpected error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select id from "user" where email = $1`)).
					WithArgs(email).
					WillReturnError(errors.New("db error"))
			},
			expected:    false,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.IsUserExist(context.Background(), email)
			assert.Equal(t, tt.expected, result)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, cleanup := mocks.SetupMockDB(t)
	defer cleanup()

	repo := &repository.PostgresUserRepository{Db: db}
	email := "abobus@mail.ru"
	loginData := models.LoginData{Email: email}
	user := models.User{
		Id:       uuid.New().String(),
		Email:    "abobus@mail.ru",
		Password: "superMegaHashUnrealNoWayReally?HashedPassword",
		Salt:     "saltySalt",
		Role:     string(models.Client),
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    models.User
		expectedErr bool
	}{
		{
			name: "user found",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "salt", "role"}).
					AddRow(user.Id, user.Email, user.Password, user.Salt, user.Role)
				mock.ExpectQuery(regexp.QuoteMeta(`select id, email, password, salt, role from "user" where email = $1`)).
					WithArgs(email).
					WillReturnRows(rows)
			},
			expected:    user,
			expectedErr: false,
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, email, password, salt, role from "user" where email = $1`)).
					WithArgs(email).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    models.User{},
			expectedErr: false,
		},
		{
			name: "db error",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select id, email, password, salt, role from "user" where email = $1`)).
					WithArgs(email).
					WillReturnError(errors.New("db error"))
			},
			expected:    models.User{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.GetUserByEmail(context.Background(), loginData)
			assert.Equal(t, tt.expected, result)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

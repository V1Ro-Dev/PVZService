package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"pvz/config/postgres"
	"pvz/internal/models"
	"pvz/pkg/logger"
)

const (
	CreateUserQuery = `
		insert into "user" (id, email, password, salt, role) values ($1, $2, $3, $4, $5)
	`

	IsUserExistQuery = `
		select id from "user" where email = $1
	`

	GetUserQuery = `
		select id, email, password, salt, role from "user" where email = $1
	`
)

type PostgresUserRepository struct {
	Db *sql.DB
}

func NewPostgresUserRepository() *PostgresUserRepository {
	db, err := sql.Open("pgx", postgres.NewPostgresConfig().GetURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return &PostgresUserRepository{Db: db}
}

func (p *PostgresUserRepository) Close() {
	p.Db.Close()
}

func (p *PostgresUserRepository) CreateUser(ctx context.Context, user models.User) error {
	_, err := p.Db.ExecContext(ctx, CreateUserQuery, user.Id, user.Email, user.Password, user.Salt, user.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return newErr
		}
		logger.Error(ctx, fmt.Sprintf("Error creating user: %s", err.Error()))
		return fmt.Errorf("unable to create user: %v", err)
	}

	return nil
}

func (p *PostgresUserRepository) IsUserExist(ctx context.Context, email string) (bool, error) {
	logger.Info(ctx, fmt.Sprintf("Checking user existance by email: %s", email))

	var userId uuid.UUID
	err := p.Db.QueryRowContext(ctx, IsUserExistQuery, email).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info(ctx, fmt.Sprintf("User with email: %s does not exist", email))
			return false, nil
		}
		logger.Error(ctx, fmt.Sprintf("unable to get user info: %v", err))
		return false, errors.New("unable to get friends info")
	}

	logger.Error(ctx, fmt.Sprintf("User with email: %s exists. His Id is %s", email, userId.String()))
	return true, nil
}

func (p *PostgresUserRepository) GetUserByEmail(ctx context.Context, logInData models.LoginData) (models.User, error) {
	logger.Info(ctx, fmt.Sprintf("Trying to get user by email: %s", logInData.Email))

	var user models.User
	err := p.Db.QueryRowContext(ctx, GetUserQuery, logInData.Email).Scan(&user.Id,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(ctx, fmt.Sprintf("User with email: %s does not exist", logInData.Email))
			return models.User{}, nil
		}
		logger.Error(ctx, fmt.Sprintf("unable to get user info: %v", err))
		return models.User{}, errors.New("unable to get friends info")
	}

	logger.Info(ctx, fmt.Sprintf("Successfully got info about user with email: %s", logInData.Email))
	return user, nil
}

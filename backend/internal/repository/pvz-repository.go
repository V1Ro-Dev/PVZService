package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"pvz/config/postgres"
	"pvz/internal/models"
	"pvz/pkg/logger"
)

const (
	CreatePvzQuery = `
		INSERT INTO pvz (id, registration_date, city) values ($1, $2, $3)
	`
)

type PostgresPvzRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresPvzRepository() *PostgresPvzRepository {
	connPool, err := pgxpool.New(context.Background(), postgres.NewPostgresConfig().GetURL())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return &PostgresPvzRepository{connPool: connPool}
}

func (p *PostgresPvzRepository) Close() {
	p.connPool.Close()
}

func (p *PostgresPvzRepository) CreatePvz(ctx context.Context, pvzData models.Pvz) error {
	logger.Info(ctx, fmt.Sprintf("Trying to create pvz with Id: %s", pvzData.Id))

	_, err := p.connPool.Exec(ctx, CreatePvzQuery, pvzData.Id, pvzData.RegistrationDate, pvzData.City)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return newErr
		}
		logger.Error(ctx, fmt.Sprintf("Error creating pvz: %s", err.Error()))
		return fmt.Errorf("unable to create pvz: %v", err)
	}
	
	return nil
}

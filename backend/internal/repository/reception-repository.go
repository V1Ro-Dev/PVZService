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
	CreateReceptiontQuery = `
		insert into reception (id, reception_datetime, pvz_id, status)
		select $1, $2, $3, $4
		where not exists (
			select id from reception
			where pvz_id = $3 AND status = $4
		)
	`
)

type PostgresReceptionRepository struct {
	connPool *pgxpool.Pool
}

func NewPostgresReceptionRepository() *PostgresReceptionRepository {
	connPool, err := pgxpool.New(context.Background(), postgres.NewPostgresConfig().GetURL())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return &PostgresReceptionRepository{connPool: connPool}
}

func (p *PostgresReceptionRepository) Close() {
	p.connPool.Close()
}

func (p *PostgresReceptionRepository) CreateReception(ctx context.Context, reception models.Reception) error {
	logger.Info(ctx, "Trying to create reception")

	commandTag, err := p.connPool.Exec(ctx, CreateReceptiontQuery, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return newErr
		}
		logger.Error(ctx, fmt.Sprintf("Error creating reception: %s", err.Error()))
		return fmt.Errorf("unable to create reception: %v", err)
	}

	if rows := commandTag.RowsAffected(); rows == 0 {
		logger.Error(ctx, "One reception was not closed or non-existing pvzId was given")
		return errors.New("one reception was not closed or non-existing pvzId was given")
	}

	logger.Info(ctx, fmt.Sprintf("Successfully created reception with Id: %s", reception.Id))
	return nil
}

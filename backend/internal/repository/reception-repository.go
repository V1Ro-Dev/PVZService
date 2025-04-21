package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"pvz/config/postgres"
	"pvz/internal/models"
	"pvz/pkg/logger"
)

const (
	CreateReceptionQuery = `
		insert into reception (id, reception_datetime, pvz_id, status)
		select $1, $2, $3, $4
		where not exists (
			select id from reception
			where pvz_id = $3 AND status = $4
		)
	`

	GetOpenReceptionQuery = `
		select id, reception_datetime, pvz_id, status
		from reception
		where pvz_id = $1 and status = $2
	`

	AddProductToOpenReceptionQuery = `
		insert into product (id, received_at, type, reception_id)
		values ($1, $2, $3, $4)
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

	commandTag, err := p.connPool.Exec(ctx, CreateReceptionQuery, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
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

func (p *PostgresReceptionRepository) GetOpenReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error) {
	logger.Info(ctx, "Trying to get open reception")

	var reception models.Reception
	err := p.connPool.QueryRow(ctx, GetOpenReceptionQuery, pvzId, models.InProgress).Scan(&reception.Id,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info(ctx, fmt.Sprintf("There is no opened reception for this pvzId: %s", pvzId.String()))
			return models.Reception{}, nil
		}
		logger.Error(ctx, fmt.Sprintf("unable to get open reception: %v", err))
		return models.Reception{}, errors.New("unable to get open reception")
	}

	logger.Info(ctx, fmt.Sprintf("Successfully got open reception with id: %s", reception.Id.String()))
	return reception, nil
}

func (p *PostgresReceptionRepository) AddProduct(ctx context.Context, product models.Product) error {
	logger.Info(ctx, "Trying to add product")

	_, err := p.connPool.Exec(ctx, AddProductToOpenReceptionQuery, product.Id, product.DateTime, product.ProductType, product.ReceptionId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return newErr
		}

		logger.Error(ctx, fmt.Sprintf("Error adding product to reception with id: %s. Error: %s", product.ReceptionId, err.Error()))
		return fmt.Errorf("unable to add product: %v", err)
	}

	logger.Info(ctx, fmt.Sprintf("Successfully added product %s to opened reception with id: %s", product.ProductType, product.ReceptionId))
	return nil
}

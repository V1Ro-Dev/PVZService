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

	DeleteLastProductFromOpenReceptionQuery = `
		delete from product
		where id = (
			select id from product
			where reception_id = $1
			order by received_at desc
			limit 1
		)
		returning id, received_at, type, reception_id
	`

	CloseReceptionQuery = `
		update reception set status = $2
		where id = $1
	`
)

type PostgresReceptionRepository struct {
	Db *sql.DB
}

func NewPostgresReceptionRepository() *PostgresReceptionRepository {
	db, err := sql.Open("pgx", postgres.NewPostgresConfig().GetURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return &PostgresReceptionRepository{Db: db}
}

func (p *PostgresReceptionRepository) Close() {
	p.Db.Close()
}

func (p *PostgresReceptionRepository) CreateReception(ctx context.Context, reception models.Reception) error {
	logger.Info(ctx, "Trying to create reception")

	commandTag, err := p.Db.ExecContext(ctx, CreateReceptionQuery, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
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

	if rows, _ := commandTag.RowsAffected(); rows == 0 {
		logger.Error(ctx, "One reception was not closed or non-existing pvzId was given")
		return errors.New("one reception was not closed or non-existing pvzId was given")
	}

	logger.Info(ctx, fmt.Sprintf("Successfully created reception with Id: %s", reception.Id))
	return nil
}

func (p *PostgresReceptionRepository) GetOpenReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error) {
	logger.Info(ctx, "Trying to get open reception")

	var reception models.Reception
	if err := p.Db.QueryRowContext(ctx, GetOpenReceptionQuery, pvzId, models.InProgress).Scan(&reception.Id,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(ctx, fmt.Sprintf("There is no opened reception for this pvzId: %s", pvzId.String()))
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

	_, err := p.Db.ExecContext(ctx, AddProductToOpenReceptionQuery, product.Id, product.DateTime, product.ProductType, product.ReceptionId)
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

func (p *PostgresReceptionRepository) RemoveProduct(ctx context.Context, receptionId uuid.UUID) error {
	logger.Info(ctx, "Trying to remove product")

	var product models.Product
	if err := p.Db.QueryRowContext(ctx, DeleteLastProductFromOpenReceptionQuery, receptionId).Scan(&product.Id,
		&product.DateTime,
		&product.ProductType,
		&product.ReceptionId,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(ctx, fmt.Sprintf("There is no active products for this receptionId: %s", receptionId.String()))
			return errors.New("there is no active products for this receptionId")
		}
		logger.Error(ctx, fmt.Sprintf("unable to delete last product for receptionId: %s. Error: %s", receptionId.String(), err.Error()))
		return errors.New("unable to delete last product")
	}

	logger.Info(ctx, fmt.Sprintf("Successfully deleted last product with params: ID: %s, ReceivedAt: %s, Type: %s, ReceptionId: %s", product.Id.String(), product.DateTime, product.ProductType, product.ReceptionId.String()))
	return nil
}

func (p *PostgresReceptionRepository) CloseReception(ctx context.Context, receptionData models.Reception) error {
	logger.Info(ctx, "Trying to close reception")

	_, err := p.Db.ExecContext(ctx, CloseReceptionQuery, receptionData.Id, models.Closed)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return newErr
		}

		logger.Error(ctx, fmt.Sprintf("Error closint reception with id: %s. Error: %s", receptionData.Id, err.Error()))
		return fmt.Errorf("unable to close reception: %v", err)
	}

	logger.Info(ctx, fmt.Sprintf("Successfully closed reception with id: %s", receptionData.Id))
	return nil
}

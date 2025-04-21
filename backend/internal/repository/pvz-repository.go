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
	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/models/postgres-models"
	"pvz/pkg/logger"
)

const (
	CreatePvzQuery = `
		insert into pvz (id, registration_date, city) values ($1, $2, $3)
	`

	GetPvzInfoQuery = `
		with paginated_pvzs as (
		  select
			pvz.id as pvz_id,
			pvz.registration_date as pvz_registration_date,
			pvz.city as pvz_city
		  from pvz
		  where pvz.registration_date between $1 AND $2
		  order by pvz.id
		  limit $3 offset $4
		)

		SELECT
		  p.pvz_id,
		  p.pvz_registration_date,
		  p.pvz_city,
		  r.id,
		  r.reception_datetime,
		  r.status,
          r.pvz_id,
		  pr.id,
		  pr.received_at,
		  pr.type,
          pr.reception_id
		from paginated_pvzs p
		left join reception r on r.pvz_id = p.pvz_id
		left join product pr on pr.reception_id = r.id
		order by p.pvz_id
	`
)

type PostgresPvzRepository struct {
	Db *sql.DB
}

func NewPostgresPvzRepository() *PostgresPvzRepository {
	db, err := sql.Open("pgx", postgres.NewPostgresConfig().GetURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return &PostgresPvzRepository{Db: db}

}

func (p *PostgresPvzRepository) Close() {
	p.Db.Close()
}

func (p *PostgresPvzRepository) CreatePvz(ctx context.Context, pvzData models.Pvz) error {
	logger.Info(ctx, fmt.Sprintf("Trying to create pvz with Id: %s", pvzData.Id))

	_, err := p.Db.ExecContext(ctx, CreatePvzQuery, pvzData.Id, pvzData.RegistrationDate, pvzData.City)
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

func (p *PostgresPvzRepository) GetPvzInfo(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error) {
	logger.Info(ctx, "Trying to get pvz info")

	rows, err := p.Db.QueryContext(ctx, GetPvzInfoQuery, form.StartDate, form.EndDate, form.Limit, (form.Page-1)*form.Limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			logger.Error(ctx, newErr.Error())
			return nil, newErr
		}
		logger.Error(ctx, fmt.Sprintf("Error getting pvz info: %s", err.Error()))
		return nil, fmt.Errorf("unable to get pvz info: %v", err)
	}
	defer rows.Close()

	pvzMap := make(map[uuid.UUID]*models.PvzInfo)

	for rows.Next() {
		var (
			pvz       postgres_models.PostgresPvz
			reception postgres_models.PostgresReception
			product   postgres_models.PostgresProduct
		)

		err = rows.Scan(
			&pvz.PvzId, &pvz.PvzRegistrationDate, &pvz.PvzCity,
			&reception.ReceptionId, &reception.ReceptionTime, &reception.ReceptionStatus, &reception.PvzId,
			&product.ProductId, &product.ProductReceivedAt, &product.ProductType, &product.ProductReceptionId,
		)

		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Scanning error: %s", err.Error()))
			return nil, err
		}

		pvzInfo, exists := pvzMap[pvz.PvzId]
		if !exists {
			pvzInfo = &models.PvzInfo{
				Pvz:        postgres_models.ToPvz(pvz),
				Receptions: []models.ReceptionProducts{},
			}
			pvzMap[pvz.PvzId] = pvzInfo
		}

		if reception.ReceptionId == uuid.Nil {
			continue
		}

		var receptionFound *models.ReceptionProducts
		for i := range pvzInfo.Receptions {
			if pvzInfo.Receptions[i].Reception.Id == reception.ReceptionId {
				receptionFound = &pvzInfo.Receptions[i]
				break
			}
		}

		if receptionFound == nil {
			newReception := models.ReceptionProducts{
				Reception: postgres_models.ToReception(reception),
				Products:  []models.Product{},
			}
			pvzInfo.Receptions = append(pvzInfo.Receptions, newReception)
			receptionFound = &pvzInfo.Receptions[len(pvzInfo.Receptions)-1]
		}

		if product.ProductId != uuid.Nil {
			receptionFound.Products = append(receptionFound.Products, postgres_models.ToProduct(product))
		}
	}

	var result []models.PvzInfo
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	logger.Info(ctx, "Successfully get pvz info")
	return result, nil
}

package postgres_models

import (
	"database/sql"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type PostgresProduct struct {
	ProductId          uuid.UUID
	ProductReceivedAt  sql.NullTime
	ProductType        sql.NullString
	ProductReceptionId uuid.UUID
}

func ToProduct(p PostgresProduct) models.Product {
	return models.Product{
		Id:          p.ProductId,
		DateTime:    p.ProductReceivedAt.Time,
		ProductType: p.ProductType.String,
		ReceptionId: p.ProductReceptionId,
	}
}

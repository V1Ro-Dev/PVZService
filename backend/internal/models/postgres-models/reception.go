package postgres_models

import (
	"database/sql"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type PostgresReception struct {
	ReceptionId     uuid.UUID
	ReceptionTime   sql.NullTime
	ReceptionStatus sql.NullString
	PvzId           uuid.UUID
}

func ToReception(p PostgresReception) models.Reception {
	return models.Reception{
		Id:       p.ReceptionId,
		DateTime: p.ReceptionTime.Time,
		PvzId:    p.PvzId,
		Status:   models.Status(p.ReceptionStatus.String),
	}
}

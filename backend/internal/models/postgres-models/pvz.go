package postgres_models

import (
	"database/sql"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type PostgresPvz struct {
	PvzId               uuid.UUID
	PvzRegistrationDate sql.NullTime
	PvzCity             sql.NullString
}

func ToPvz(p PostgresPvz) models.Pvz {
	return models.Pvz{
		Id:               p.PvzId,
		RegistrationDate: p.PvzRegistrationDate.Time,
		City:             p.PvzCity.String,
	}
}

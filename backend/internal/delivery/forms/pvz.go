package forms

import (
	"pvz/internal/models"
	"time"

	"github.com/google/uuid"
)

type PvzForm struct {
	Id               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

func ToPvzForm(pvz models.Pvz) PvzForm {
	return PvzForm{
		Id:               pvz.Id,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}

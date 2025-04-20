package forms

import (
	"time"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type PvzForm struct {
	Id               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

func ToPvzForm(pvz models.Pvz) PvzForm {
	return PvzForm{
		Id:               pvz.Id,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}

package forms

import (
	"time"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type ReceptionForm struct {
	PvzId uuid.UUID `json:"pvzId"`
}

type ReceptionFormOut struct {
	Id       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
}

func ToReceptionFormOut(reception models.Reception) ReceptionFormOut {
	return ReceptionFormOut{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzId,
		Status:   string(reception.Status),
	}
}

type ReceptionProductsFormOut struct {
	Reception ReceptionFormOut `json:"reception"`
	Products  []ProductFormOut `json:"products"`
}

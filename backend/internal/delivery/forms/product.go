package forms

import (
	"time"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type ProductForm struct {
	PvzId uuid.UUID `json:"pvzId"`
	Type  string    `json:"type"`
}

type ProductFormOut struct {
	Id          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	ProductType string    `json:"productType"`
	ReceptionId uuid.UUID `json:"receptionId"`
}

func ToProductFormOut(product models.Product) ProductFormOut {
	return ProductFormOut{
		Id:          product.Id,
		DateTime:    product.DateTime,
		ProductType: product.ProductType,
		ReceptionId: product.ReceptionId,
	}
}

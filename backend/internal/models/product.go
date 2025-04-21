package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	DateTime    time.Time
	ProductType string
	ReceptionId uuid.UUID
}

package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id          uuid.UUID
	DateTime    time.Time
	ProductType string
	ReceptionId uuid.UUID
}

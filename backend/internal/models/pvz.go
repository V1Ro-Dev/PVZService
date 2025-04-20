package models

import (
	"time"

	"github.com/google/uuid"
)

type Pvz struct {
	Id               uuid.UUID
	RegistrationDate time.Time
	City             string
}

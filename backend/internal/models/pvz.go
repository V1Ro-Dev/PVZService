package models

import (
	"github.com/google/uuid"
	"time"
)

type Pvz struct {
	Id               uuid.UUID
	RegistrationDate time.Time
	City             string
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	Moderator Role = "moderator"
	Employee  Role = "employee"
	Client    Role = "client"
)

type Pvz struct {
	Id               uuid.UUID
	RegistrationDate time.Time
	City             string
}

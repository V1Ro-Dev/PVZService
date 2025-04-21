package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	InProgress Status = "in_progress"
	Closed     Status = "close"
)

type Reception struct {
	Id       uuid.UUID
	DateTime time.Time
	PvzId    uuid.UUID
	Status   Status
}

type ReceptionProducts struct {
	Reception Reception
	Products  []Product
}

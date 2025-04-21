package models

import (
	"github.com/google/uuid"
	"time"
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

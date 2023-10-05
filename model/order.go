package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id         uuid.UUID
	Status     string
	ClientId   uuid.UUID
	Creatde_at *time.Time
}

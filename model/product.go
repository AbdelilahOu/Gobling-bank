package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	Tva         float64
	Created_at  *time.Time
}

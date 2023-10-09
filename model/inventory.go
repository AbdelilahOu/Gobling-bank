package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryMvm struct {
	Id        uuid.UUID
	Quantity  int
	CreatedAt *time.Time
	ProductId uuid.UUID
}

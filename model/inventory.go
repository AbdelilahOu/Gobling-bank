package model

import (
	"github.com/google/uuid"
)

type InventoryMvm struct {
	Id        uuid.UUID `db:"id"`
	Quantity  int       `db:"quantity"`
	CreatedAt string    `db:"created_at"`
	ProductId uuid.UUID `db:"product_id"`
}

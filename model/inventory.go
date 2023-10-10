package model

import (
	"github.com/google/uuid"
)

type InventoryMvm struct {
	Id        uuid.UUID `sql:"id"`
	Quantity  int       `sql:"quantity"`
	CreatedAt string    `sql:"created_at"`
	ProductId uuid.UUID `sql:"product_id"`
}

package model

import (
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Tva         float64   `db:"tva"`
	CreatedAt   string    `db:"created_at"`
}

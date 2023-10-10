package model

import (
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `sql:"id"`
	Name        string    `sql:"name"`
	Description string    `sql:"description"`
	Price       float64   `sql:"price"`
	Tva         float64   `sql:"tva"`
	CreatedAt   string    `sql:"created_at"`
}

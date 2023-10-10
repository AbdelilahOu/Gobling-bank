package model

import "github.com/google/uuid"

type OrderItem struct {
	Id          uuid.UUID `db:"id"`
	ProductId   uuid.UUID `db:"product_id"`
	NewPrice    float64   `db:"new_price"`
	Quantity    int       `db:"quantity"`
	OrderId     uuid.UUID `db:"order_id"`
	InventoryId uuid.UUID `db:"inventory_id"`
}

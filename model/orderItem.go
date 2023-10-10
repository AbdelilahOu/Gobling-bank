package model

import "github.com/google/uuid"

type OrderItem struct {
	Id          uuid.UUID `sql:"id"`
	ProductId   uuid.UUID `sql:"product_id"`
	NewPrice    float64   `sql:"new_price"`
	Quantity    int       `sql:"quantity"`
	OrderId     uuid.UUID `sql:"order_id"`
	InventoryId uuid.UUID `sql:"inventory_id"`
}

package model

import "github.com/google/uuid"

type OrderItem struct {
	Id          uuid.UUID
	ProductId   uuid.UUID
	NewPrice    float64
	Quantity    int
	OrderId     uuid.UUID
	InventoryId uuid.UUID
}

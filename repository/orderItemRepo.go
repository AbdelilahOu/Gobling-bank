package repository

import (
	"context"
	"database/sql"

	"github.com/AbdelilahOu/GoThingy/model"
)

type OrderItemRepo struct {
	DB *sql.DB
}

func (repo *OrderItemRepo) Insert(ctx context.Context, orderItem model.OrderItem) error {
	return nil
}

func (repo *OrderItemRepo) Update(ctx context.Context, orderItem model.OrderItem, id string) error {
	return nil
}

func (repo *OrderItemRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *OrderItemRepo) Select(ctx context.Context, id string) error {
	return nil
}

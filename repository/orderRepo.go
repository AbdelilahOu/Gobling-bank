package repository

import (
	"context"
	"database/sql"

	"github.com/AbdelilahOu/GoThingy/model"
)

type OrderRepo struct {
	DB *sql.DB
}

func (repo *OrderRepo) Insert(ctx context.Context, order model.Order) error {
	repo.DB.Query("INSERT INTO orders")

	return nil
}

func (repo *OrderRepo) Update(ctx context.Context, order model.Order, id string) error {
	return nil
}

func (repo *OrderRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *OrderRepo) Select(ctx context.Context, id string) error {
	return nil
}

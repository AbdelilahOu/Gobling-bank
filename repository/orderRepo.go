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
	return nil
}

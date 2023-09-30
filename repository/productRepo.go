package repository

import (
	"context"
	"database/sql"

	"github.com/AbdelilahOu/GoThingy/model"
)

type ProductRepo struct {
	DB *sql.DB
}

func (repo *ProductRepo) Insert(ctx context.Context, product model.Product) error {
	return nil
}

func (repo *ProductRepo) Update(ctx context.Context, product model.Product, id string) error {
	return nil
}

func (repo *ProductRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *ProductRepo) Select(ctx context.Context, id string) error {
	return nil
}

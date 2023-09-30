package repository

import (
	"context"
	"database/sql"

	"github.com/AbdelilahOu/GoThingy/model"
)

type InventoryRepo struct {
	DB *sql.DB
}

func (repo *InventoryRepo) Insert(ctx context.Context, inventory model.InventoryMvm) error {
	return nil
}

func (repo *InventoryRepo) Update(ctx context.Context, inventory model.InventoryMvm, id string) error {
	return nil
}

func (repo *InventoryRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *InventoryRepo) Select(ctx context.Context, id string) error {
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryRepo struct {
	DB *sqlx.DB
}

type GetIAllResult struct {
	Inventory []model.InventoryMvm
	Cursor    uint64
}

func (repo *InventoryRepo) Insert(ctx context.Context, inventory model.InventoryMvm) (uuid.UUID, error) {
	_, err := repo.DB.Exec("INSERT INTO inventory_mouvements (id,product_id,quantity) VALUES ($1,$2,$3)", inventory.Id, inventory.ProductId, inventory.Quantity)
	if err != nil {
		fmt.Println("error inserting inventory", err)
		return uuid.UUID{}, err
	}
	return inventory.Id, nil
}

func (repo *InventoryRepo) Update(ctx context.Context, inventory model.InventoryMvm, id string) error {
	_, err := repo.DB.Exec("UPDATE inventory_mouvements SET quantity=$1 WHERE id=$2", inventory.Quantity, id)

	if err != nil {
		fmt.Println("error updating inventory :", err)
		return err
	}
	return nil
}

func (repo *InventoryRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.DB.Exec("DELETE FROM inventory_mouvements WHERE id = $1", id)
	if err != nil {
		fmt.Println("error deleting inventory", err)
		return err
	}
	return nil
}

func (repo *InventoryRepo) Select(ctx context.Context, id string) (model.InventoryMvm, error) {
	// execute
	var inventory model.InventoryMvm
	err := repo.DB.Select(inventory, "SELECT * FROM inventory_mouvements WHERE id = $1", id)
	// var
	if err == sql.ErrNoRows {
		fmt.Println("no redcord exisist", err)
		return model.InventoryMvm{}, errorMessages.RecordDoesntExist

	}
	//
	if err != nil {
		fmt.Println("error scanning inventory", err)
		return model.InventoryMvm{}, err
	}
	//
	return inventory, nil
}

func (repo *InventoryRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetIAllResult, error) {
	// get inventory_mouvements
	var inventories []model.InventoryMvm
	err := repo.DB.Select(&inventories, "SELECT * FROM inventory_mouvements WHERE id > $1 LIMIT $2", cursor, size)
	if err != nil {
		fmt.Println("error getting inventory_mouvements", err)
		return GetIAllResult{}, err
	}
	// last result
	return GetIAllResult{
		Inventory: inventories,
	}, nil
}

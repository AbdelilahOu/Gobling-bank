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
	row := repo.DB.QueryRow("SELECT * FROM inventory_mouvements WHERE id = $1", id)
	// var
	var c model.InventoryMvm
	// get inventory
	err := row.Scan()
	// check for err
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
	return c, nil
}

func (repo *InventoryRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetIAllResult, error) {
	// get inventory_mouvements
	rows, err := repo.DB.Query("SELECT * FROM inventory_mouvements WHERE id > $1 LIMIT $2", cursor, size)
	if err != nil {
		fmt.Println("error getting inventory_mouvements", err)
		return GetIAllResult{}, err
	}
	// close rows after
	defer rows.Close()
	// get inventory_mouvements as model.InventoryMvm
	var Inventory []model.InventoryMvm
	for rows.Next() {
		var c model.InventoryMvm
		// scane
		err := rows.Scan()
		if err != nil {
			fmt.Println("error scanning Inventory", err)
			return GetIAllResult{}, err
		}
		//
		Inventory = append(Inventory, c)

	}
	// error while eterating
	err = rows.Err()
	if err != nil {
		fmt.Println("error eterating over rows")
	}
	// last result
	return GetIAllResult{
		Inventory: Inventory,
	}, nil
}

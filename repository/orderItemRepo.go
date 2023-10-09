package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
)

type OrderItemRepo struct {
	DB *sql.DB
}

func (repo *OrderItemRepo) Insert(ctx context.Context, orderItem model.OrderItem) error {
	_, err := repo.DB.Exec("INSERT INTO order_items (id,order_id, product_id, quantity, new_price, inventory_id) VALUES ($1, $2, $3, $4, $5, $6)", orderItem.Id, orderItem.OrderId, orderItem.ProductId, orderItem.Quantity, orderItem.NewPrice, orderItem.InventoryId)
	if err != nil {
		fmt.Println("error inserting orderItem", err)
		return err
	}
	return nil
}

func (repo *OrderItemRepo) Update(ctx context.Context, orderItem model.OrderItem, id string) error {
	_, err := repo.DB.Exec("UPDATE order_items SET quantity = $1, new_price = $2 WHERE id = $3", orderItem.Quantity, orderItem.NewPrice)

	if err != nil {
		fmt.Println("error updating orderItem :", err)
		return err
	}
	return nil
}

func (repo *OrderItemRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.DB.Exec("DELETE FROM order_items WHERE id = $1", id)
	if err != nil {
		fmt.Println("error deleting orderItem", err)
		return err
	}
	return nil
}

func (repo *OrderItemRepo) Select(ctx context.Context, id string) (model.OrderItem, error) {
	// execute
	row := repo.DB.QueryRow("SELECT * FROM order_items WHERE id = $1", id)
	// var
	var c model.OrderItem
	// get orderItem
	err := row.Scan()
	// check for err
	if err == sql.ErrNoRows {
		fmt.Println("no redcord exisist", err)
		return model.OrderItem{}, errorMessages.RecordDoesntExist

	}
	//
	if err != nil {
		fmt.Println("error scanning orderItem", err)
		return model.OrderItem{}, err
	}
	//
	return c, nil
}

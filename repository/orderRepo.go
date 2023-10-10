package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	DB *sqlx.DB
}

type GetOAllResult struct {
	Orders []model.Order
	Cursor uint64
}

func (repo *OrderRepo) Insert(ctx context.Context, order model.Order) error {
	_, err := repo.DB.Exec("INSERT INTO orders (id, client_id, status) VALUES ($1, $2, $3)", order.Id, order.ClientId, order.Status)
	if err != nil {
		fmt.Println("error inserting order", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Update(ctx context.Context, order model.Order, id string) error {
	_, err := repo.DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", order.Status, id)

	if err != nil {
		fmt.Println("error updating order :", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.DB.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		fmt.Println("error deleting order", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Select(ctx context.Context, id string) (model.Order, error) {
	// execute
	var order model.Order
	err := repo.DB.Select(&order, "SELECT * FROM orders WHERE id = $1", id)
	// var
	if err == sql.ErrNoRows {
		fmt.Println("no redcord exisist", err)
		return model.Order{}, errorMessages.RecordDoesntExist

	}
	//
	if err != nil {
		fmt.Println("error scanning order", err)
		return model.Order{}, err
	}
	//
	return order, nil
}

func (repo *OrderRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetOAllResult, error) {
	// get orders
	var orders []model.Order
	err := repo.DB.Select(&orders, "SELECT id, client_id, status, created_at  FROM orders WHERE id > $1 LIMIT $2", cursor, size)
	if err != nil {
		fmt.Println("error getting orders", err)
		return GetOAllResult{}, err
	}
	// ok
	return GetOAllResult{
		Orders: orders,
	}, nil
}

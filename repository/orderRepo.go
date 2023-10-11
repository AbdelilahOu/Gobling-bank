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
	InsertQuery := "INSERT INTO orders (id, client_id, status) VALUES ($1, $2, $3)"
	_, err := repo.DB.Exec(InsertQuery, order.Id, order.ClientId, order.Status)
	if err != nil {
		fmt.Println("error inserting order", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Update(ctx context.Context, order model.Order, id string) error {
	UpdateQuery := "UPDATE orders SET status = $1 WHERE id = $2"
	_, err := repo.DB.Exec(UpdateQuery, order.Status, id)

	if err != nil {
		fmt.Println("error updating order :", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Delete(ctx context.Context, id string) error {
	DeleteQuery := "DELETE FROM orders WHERE id = $1"
	_, err := repo.DB.Exec(DeleteQuery, id)
	if err != nil {
		fmt.Println("error deleting order", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) Select(ctx context.Context, id string) (model.Order, error) {
	// execute
	SelectQuery := "SELECT * FROM orders WHERE id = $1"
	var order model.Order
	err := repo.DB.Select(&order, SelectQuery, id)
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
	SelectAllQuery := "SELECT id, client_id, status, created_at  FROM orders WHERE id > $1 LIMIT $2"
	var orders []model.Order
	err := repo.DB.Select(&orders, SelectAllQuery, cursor, size)
	if err != nil {
		fmt.Println("error getting orders", err)
		return GetOAllResult{}, err
	}
	// ok
	return GetOAllResult{
		Orders: orders,
	}, nil
}

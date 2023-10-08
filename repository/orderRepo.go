package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
)

type OrderRepo struct {
	DB *sql.DB
}

type GetOAllResult struct {
	Orders []model.Order
	Cursor uint64
}

func (repo *OrderRepo) Insert(ctx context.Context, order model.Order) error {
	_, err := repo.DB.Exec("INSERT INTO orders (id, client_id, status) VALUES ($1, $2, $3, $4)", order.Id, order.ClientId, order.Status)
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
	row := repo.DB.QueryRow("SELECT * FROM orders WHERE id = $1", id)
	// var
	var c model.Order
	// get order
	err := row.Scan()
	// check for err
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
	return c, nil
}

func (repo *OrderRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetOAllResult, error) {
	// get orders
	rows, err := repo.DB.Query("SELECT * FROM orders WHERE id > $1 LIMIT $2", cursor, size)
	if err != nil {
		fmt.Println("error getting orders", err)
		return GetOAllResult{}, err
	}
	// close rows after
	defer rows.Close()
	// get orders as model.Order
	var orders []model.Order
	for rows.Next() {
		var c model.Order
		// scane
		err := rows.Scan()
		if err != nil {
			fmt.Println("error scanning orders", err)
			return GetOAllResult{}, err
		}
		//
		orders = append(orders, c)

	}
	// error while eterating
	err = rows.Err()
	if err != nil {
		fmt.Println("error eterating over rows")
	}
	// last result
	return GetOAllResult{
		Orders: orders,
	}, nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AbdelilahOu/GoThingy/model"
)

type ClientRepo struct {
	DB *sql.DB
}

func (repo *ClientRepo) Insert(ctx context.Context, client model.Client) error {
	_, err := repo.DB.Exec("INSERT INTO clients (id, firstname, lastname, email, phone) VALUES ($1, $2, $3, $4, $5)", client.Id, client.Firstname, client.Lastname, client.Email, client.Phone)
	if err != nil {
		fmt.Println("error inserting client", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Update(ctx context.Context, client model.Client, id string) error {
	_, err := repo.DB.Exec("UPDATE clients SET firstname = $1, lastname = $2, email = $3, phone = $4 WHERE id = $5", client.Firstname, client.Lastname, client.Email, client.Phone, id)
	if err != nil {
		fmt.Println("error updating client :", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.DB.Exec("DELETE FROM clients WHERE id = $1", id)
	if err != nil {
		fmt.Println("error deleting client", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Select(ctx context.Context, id string) (model.Client, error) {
	// execute
	row := repo.DB.QueryRow("SELECT firstname, lastname, email, phone, created_at, adress  FROM clients WHERE id = $1", id)
	// var
	var c model.Client
	// get client
	err := row.Scan(&c.Id, &c.Firstname, &c.Lastname, &c.Email, &c.Phone, &c.Created_at, &c.Adress)
	// check for err
	if err != nil {
		fmt.Println("error scanning client", err)
		return model.Client{}, err
	}
	//
	return c, nil
}

type GetCAllResult struct {
	Clients []model.Client
	Cursor  uint64
}

func (repo *ClientRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetCAllResult, error) {
	// get clients
	rows, err := repo.DB.Query("SELECT * FROM clients WHERE id > $1", cursor)
	if err != nil {
		fmt.Println("error getting clients", err)
		return GetCAllResult{}, err
	}
	// close rows after
	defer rows.Close()
	// get clients as model.client
	var clients []model.Client
	for rows.Next() {
		var c model.Client
		// scane
		err := rows.Scan(&c.Id, &c.Firstname, &c.Lastname, &c.Email, &c.Phone)
		if err != nil {
			fmt.Println("error scanning clients", err)
			return GetCAllResult{}, err
		}
		//
		clients = append(clients, c)

	}
	// error while eterating
	err = rows.Err()
	if err != nil {
		fmt.Println("error eterating over rows")
	}
	// last result
	return GetCAllResult{
		Clients: clients,
	}, nil
}

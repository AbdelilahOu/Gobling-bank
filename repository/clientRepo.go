package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/jmoiron/sqlx"
)

type ClientRepo struct {
	DB *sqlx.DB
}

type GetCAllResult struct {
	Clients []model.Client
	Cursor  uint64
}

func (repo *ClientRepo) Insert(ctx context.Context, client model.Client) error {
	InsertQuery := "INSERT INTO clients (id, firstname, lastname, email, phone) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.DB.Exec(InsertQuery, client.Id, client.Firstname, client.Lastname, client.Email, client.Phone)
	if err != nil {
		fmt.Println("error inserting client", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Update(ctx context.Context, client model.Client, id string) error {
	UpdateQuery := "UPDATE clients SET firstname = $1, lastname = $2, email = $3, phone = $4 WHERE id = $5"
	_, err := repo.DB.Exec(UpdateQuery, client.Firstname, client.Lastname, client.Email, client.Phone, id)

	if err != nil {
		fmt.Println("error updating client :", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Delete(ctx context.Context, id string) error {
	DeleteQuery := "DELETE FROM clients WHERE id = $1"
	_, err := repo.DB.Exec(DeleteQuery, id)
	if err != nil {
		fmt.Println("error deleting client", err)
		return err
	}
	return nil
}

func (repo *ClientRepo) Select(ctx context.Context, id string) (model.Client, error) {
	SelectQuery := "SELECT firstname, lastname, email, phone, created_at, address  FROM clients WHERE id = $1"
	// execute
	var c model.Client
	err := repo.DB.Select(&c, SelectQuery, id)
	// check for err
	if err == sql.ErrNoRows {
		fmt.Println("no redcord exisist", err)
		return model.Client{}, errorMessages.RecordDoesntExist

	}
	//
	if err != nil {
		fmt.Println("error scanning client", err)
		return model.Client{}, err
	}
	//
	return c, nil
}

func (repo *ClientRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetCAllResult, error) {
	SelectAllQuery := "SELECT * FROM clients WHERE id > $1 LIMIT $2"
	// get clients
	var clients []model.Client
	err := repo.DB.Select(&clients, SelectAllQuery, cursor, size)
	if err != nil {
		fmt.Println("error getting clients", err)
		return GetCAllResult{}, err
	}
	// last result
	return GetCAllResult{
		Clients: clients,
	}, nil
}

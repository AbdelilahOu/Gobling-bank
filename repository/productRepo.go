package repository

import (
	"context"
	"database/sql"
	"fmt"

	errorMessages "github.com/AbdelilahOu/GoThingy/constants"
	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/jmoiron/sqlx"
)

type ProductRepo struct {
	DB *sqlx.DB
}

type GetPAllResult struct {
	Products []model.Product
	Cursor   uint64
}

func (repo *ProductRepo) Insert(ctx context.Context, product model.Product) error {
	InsertQuery := "INSERT INTO products (id, name, description, price, tva) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.DB.Exec(InsertQuery, product.Id, product.Name, product.Description, product.Price, product.Tva)
	if err != nil {
		fmt.Println("error inserting product", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Update(ctx context.Context, product model.Product, id string) error {
	UpdateQuery := "UPDATE products SET name = $1, description = $2, price = $3, tva = $4 WHERE id = $5"
	_, err := repo.DB.Exec(UpdateQuery, product.Name, product.Description, product.Price, product.Tva, id)
	// check for error
	if err != nil {
		fmt.Println("error updating product :", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Delete(ctx context.Context, id string) error {
	DeleteQuery := "DELETE FROM products WHERE id = $1"
	_, err := repo.DB.Exec(DeleteQuery, id)
	if err != nil {
		fmt.Println("error deleting product", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Select(ctx context.Context, id string) (model.Product, error) {
	// execute
	SelectQuery := "SELECT * FROM products WHERE id = $1"
	var product model.Product
	err := repo.DB.Select(&product, SelectQuery, id)
	//
	if err == sql.ErrNoRows {
		fmt.Println("no redcord exisist", err)
		return model.Product{}, errorMessages.RecordDoesntExist

	}
	//
	if err != nil {
		fmt.Println("error scanning product", err)
		return model.Product{}, err
	}
	//
	return product, nil
}

func (repo *ProductRepo) SelectAll(ctx context.Context, cursor uint64, size uint64) (GetPAllResult, error) {
	// get products
	SelectAllQuery := "SELECT * FROM products WHERE id > $1 LIMIT $2"
	var products []model.Product
	err := repo.DB.Select(products, SelectAllQuery, cursor, size)
	if err != nil {
		fmt.Println("error getting products", err)
		return GetPAllResult{}, err
	}
	// last result
	return GetPAllResult{
		Products: products,
	}, nil
}

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
	_, err := repo.DB.Exec("INSERT INTO products (id, name, description, price, tva) VALUES ($1, $2, $3, $4, $5)", product.Id, product.Name, product.Description, product.Price, product.Tva)
	if err != nil {
		fmt.Println("error inserting product", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Update(ctx context.Context, product model.Product, id string) error {
	_, err := repo.DB.Exec("UPDATE products SET name = $1, description = $2, price = $3, tva = $4 WHERE id = $5", product.Name, product.Description, product.Price, product.Tva, id)
	// check for error
	if err != nil {
		fmt.Println("error updating product :", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Delete(ctx context.Context, id string) error {
	_, err := repo.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		fmt.Println("error deleting product", err)
		return err
	}
	return nil
}

func (repo *ProductRepo) Select(ctx context.Context, id string) (model.Product, error) {
	// execute
	var product model.Product
	err := repo.DB.Select(&product, "SELECT * FROM products WHERE id = $1", id)
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
	var products []model.Product
	err := repo.DB.Select(products, "SELECT * FROM products WHERE id > $1 LIMIT $2", cursor, size)
	if err != nil {
		fmt.Println("error getting products", err)
		return GetPAllResult{}, err
	}
	// last result
	return GetPAllResult{
		Products: products,
	}, nil
}

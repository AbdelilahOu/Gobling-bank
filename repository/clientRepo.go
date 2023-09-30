package repository

import (
	"context"
	"database/sql"

	"github.com/AbdelilahOu/GoThingy/model"
)

type ClientRepo struct {
	DB *sql.DB
}

func (repo *ClientRepo) Insert(ctx context.Context, client model.Client) error {
	return nil
}

func (repo *ClientRepo) Update(ctx context.Context, client model.Client, id string) error {
	return nil
}

func (repo *ClientRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *ClientRepo) Select(ctx context.Context, id string) error {
	return nil
}

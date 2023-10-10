package model

import (
	"github.com/google/uuid"
)

type Client struct {
	Id        uuid.UUID `db:"id"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Phone     string    `db:"phone"`
	Email     string    `db:"email"`
	Address   string    `db:"address"`
	CreatedAt string    `db:"created_at"`
}

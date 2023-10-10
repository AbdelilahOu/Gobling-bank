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

// type Client struct {
// 	Id        uuid.UUID  `db:"clients.id" db:"id"`
// 	Firstname string     `db:"clients.firstname" db:"firstname"`
// 	Lastname  string     `db:"clients.lastname" db:"lastname"`
// 	Phone     string     `db:"clients.phone" db:"phone"`
// 	Email     string     `db:"clients.email" db:"email"`
// 	Address    string     `db:"clients.address" db:"address"`
// 	CreatedAt string `db:"clients.created_at" db:"created_at"`
// }

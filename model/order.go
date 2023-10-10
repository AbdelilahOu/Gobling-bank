package model

import (
	"github.com/google/uuid"
)

type Order struct {
	Id        uuid.UUID `db:"id"`
	Status    string    `db:"status"`
	ClientId  uuid.UUID `db:"client_id"`
	CreatedAt string    `db:"created_at"`
}

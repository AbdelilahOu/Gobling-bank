package model

import (
	"github.com/google/uuid"
)

type Order struct {
	Id        uuid.UUID `sql:"id"`
	Status    string    `sql:"status"`
	ClientId  uuid.UUID `sql:"client_id"`
	CreatedAt string    `sql:"created_at"`
}

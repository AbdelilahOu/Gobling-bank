package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Id         uuid.UUID
	Firstname  string
	Lastname   string
	Phone      string
	Email      string
	Created_at *time.Time
}

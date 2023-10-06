package handler

import (
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Order struct {
	Repo *repository.OrderRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

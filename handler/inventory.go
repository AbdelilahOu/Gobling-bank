package handler

import (
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Inventory struct {
	Repo *repository.InventoryRepo
}

func (o *Inventory) Create(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

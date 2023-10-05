package handler

import (
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Inventory struct {
	Repo *repository.InventoryRepo
}

func (o *Inventory) Create(w http.ResponseWriter, r *http.Request) {

}

func (o *Inventory) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Inventorys")
}

func (o *Inventory) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Inventory")
}

func (o *Inventory) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update Inventory")
}

func (o *Inventory) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete Inventory")
}

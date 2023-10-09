package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/AbdelilahOu/GoThingy/repository"
	"github.com/google/uuid"
)

type Inventory struct {
	Repo *repository.InventoryRepo
}

func (o *Inventory) Create(w http.ResponseWriter, r *http.Request) {
	// body struct
	var body struct {
		Quantity  int    `json:"quantity"`
		ProductId string `json:"product_id"`
	}
	// get body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// create uuid
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Prase id
	ProductId, err := uuid.Parse(body.ProductId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create inventory
	inv, err := o.Repo.Insert(r.Context(), model.InventoryMvm{
		Id:        id,
		Quantity:  body.Quantity,
		ProductId: ProductId,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// return inventory
	res, err := json.Marshal(model.InventoryMvm{
		Id:        inv,
		Quantity:  body.Quantity,
		ProductId: ProductId,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Inventory) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Inventory) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

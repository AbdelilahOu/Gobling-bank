package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/AbdelilahOu/GoThingy/repository"
	"github.com/google/uuid"
)

type OrderItem struct {
	Repo *repository.OrderItemRepo
}

func (o *OrderItem) Create(w http.ResponseWriter, r *http.Request) {
	// bosy struct
	var body struct {
		ProductId   string  `json:"product_id"`
		OrderId     string  `json:"order_id"`
		InventoryId string  `json:"inventory_id"`
		NewPrice    float64 `json:"new_price"`
		Quantity    int     `json:"quantity"`
	}
	// get body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// uuid
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	parsedProductId, err := uuid.Parse(body.ProductId)
	if err != nil {
		fmt.Println("error parsing product id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	parsedOrderId, err := uuid.Parse(body.OrderId)
	if err != nil {
		fmt.Println("error parsing product id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	parsedInventoryId, err := uuid.Parse(body.InventoryId)
	if err != nil {
		fmt.Println("error parsing product id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create product
	product := model.OrderItem{
		Id:          id,
		ProductId:   parsedProductId,
		OrderId:     parsedOrderId,
		InventoryId: parsedInventoryId,
		NewPrice:    body.NewPrice,
		Quantity:    body.Quantity,
	}
	// add to db
	err = o.Repo.Insert(r.Context(), product)
	if err != nil {
		fmt.Println("error inserting product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	res, err := json.Marshal(product)
	if err != nil {
		fmt.Println("error marshaling product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *OrderItem) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *OrderItem) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *OrderItem) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *OrderItem) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

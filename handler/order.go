package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/AbdelilahOu/GoThingy/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Order struct {
	Repo          *repository.OrderRepo
	ItemsRepo     *repository.OrderItemRepo
	InventoryRepo *repository.InventoryRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	// body struct
	var body struct {
		ClientId string `json:"client_id"`
		Status   string `json:"status"`
		Items    []struct {
			ProductId string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			NewPrice  float64 `json:"new_price"`
		} `json:"items"`
	}
	// decode body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// create order
	now := time.Now().UTC()
	orderId, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// parse uuid
	parsedClientId, err := uuid.Parse(body.ClientId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create order
	err = o.Repo.Insert(r.Context(), model.Order{
		Id:        orderId,
		Status:    body.Status,
		ClientId:  parsedClientId,
		CreatedAt: &now,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create order items
	for _, item := range body.Items {
		inventoryId, err := uuid.NewUUID()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//
		parsedProductId, err := uuid.Parse(item.ProductId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//
		inventoryId, err = o.InventoryRepo.Insert(r.Context(), model.InventoryMvm{
			Id:        inventoryId,
			Quantity:  item.Quantity,
			CreatedAt: &now,
			ProductId: parsedProductId,
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//
		orderItemId, err := uuid.NewUUID()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//
		o.ItemsRepo.Insert(r.Context(), model.OrderItem{
			Id:          orderItemId,
			OrderId:     orderId,
			ProductId:   parsedProductId,
			NewPrice:    item.NewPrice,
			Quantity:    item.Quantity,
			InventoryId: inventoryId,
		})
	}
	w.Write([]byte(orderId.String()))
	w.WriteHeader(http.StatusCreated)
}

func (o *Order) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// get params
	idParam := chi.URLParam(r, "id")
	// check if param exist
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete order
	err := o.Repo.Delete(r.Context(), idParam)
	// check for errors
	if err != nil {
		fmt.Println("error deleting client", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.WriteHeader(http.StatusOK)
}

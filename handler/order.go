package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		CreatedAt: now.String(),
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
			CreatedAt: now.String(),
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
	// pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	// check
	if limitStr == "" {
		limitStr = "10"
	}
	// inisialise pageStr if not provided
	if pageStr == "" {
		pageStr = "0"
	}
	//convert page into int
	const decimal = 10
	const bitSize = 64
	page, err := strconv.ParseUint(pageStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := strconv.ParseUint(limitStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//
	result, err := o.Repo.SelectAll(r.Context(), page, size)
	if err != nil {
		fmt.Println("coudnt get orders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("couldnt marshal orders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	// check if param exist
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get order
	result, err := o.Repo.Select(r.Context(), idParam)
	if err != nil {
		fmt.Println("couldnt get order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// marhsal json
	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("couldnt marshal order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// body struct
	var body struct {
		Status    string `json:"status"`
		OrderItem []struct {
			Id          string  `json:"id,omitempty"`
			Quantity    int     `json:"quantity"`
			ProductId   string  `json:"product_id"`
			NewPrice    float64 `json:"new_price"`
			InventoryId string  `json:"inventory_id,omitempty"`
		}
	}
	// populat nody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get params
	idParam := chi.URLParam(r, "id")
	// check if param exist
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if client exists
	// get client
	order, err := o.Repo.Select(r.Context(), idParam)
	// check for erros
	if err != nil {
		fmt.Println("error fetching order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// update order
	err = o.Repo.Update(r.Context(), model.Order{Status: body.Status}, idParam)
	if err != nil {
		fmt.Println("error updating order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// update order items
	for _, item := range body.OrderItem {
		if item.Id == "" {
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
			parsedOrderId, err := uuid.Parse(idParam)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//
			o.ItemsRepo.Insert(r.Context(), model.OrderItem{
				Id:          orderItemId,
				OrderId:     parsedOrderId,
				ProductId:   parsedProductId,
				NewPrice:    item.NewPrice,
				Quantity:    item.Quantity,
				InventoryId: inventoryId,
			})
		} else {
			err := o.InventoryRepo.Update(r.Context(), model.InventoryMvm{Quantity: item.Quantity}, item.InventoryId)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = o.ItemsRepo.Update(r.Context(), model.OrderItem{Quantity: item.Quantity, NewPrice: item.NewPrice}, item.Id)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	//
	order.Status = body.Status
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("error marshalling order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.Write(res)
	w.WriteHeader(http.StatusOK)
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

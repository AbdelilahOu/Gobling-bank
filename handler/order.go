package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Order struct {
	Repo *repository.OrderRepo
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
	fmt.Println(body.ClientId)
}

func (o *Order) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

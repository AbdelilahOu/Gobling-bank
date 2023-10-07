package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/AbdelilahOu/GoThingy/repository"
	"github.com/google/uuid"
)

type Product struct {
	Repo *repository.ProductRepo
}

func (o *Product) Create(w http.ResponseWriter, r *http.Request) {
	// bosy struct
	var body struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Tva         float64 `json:"tva"`
	}
	// get body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// timr
	now := time.Now().UTC()
	// uuid
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create product
	product := model.Product{
		Id:          id,
		Name:        body.Name,
		Price:       body.Price,
		Description: body.Description,
		Tva:         body.Tva,
		Created_at:  &now,
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
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

func (o *Product) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (o *Product) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Product) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Product) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

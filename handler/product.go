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
	// get products
	products, err := o.Repo.SelectAll(r.Context(), page, size)
	if err != nil {
		fmt.Println("error getting products")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get json
	res, err := json.Marshal(products)
	if err != nil {
		fmt.Println("error marshaling products")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (o *Product) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// check for id
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client, err := o.Repo.Select(r.Context(), id)
	if err != nil {
		fmt.Println("error getting product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get json
	res, err := json.Marshal(client)
	if err != nil {
		fmt.Println("error marshaling product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (o *Product) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Product) DeleteByID(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error deleting product", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ok
	w.WriteHeader(http.StatusOK)
}

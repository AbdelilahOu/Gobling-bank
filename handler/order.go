package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("create order")
}

func (o *Order) GetAll(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("get orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("get order")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("update order")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("delete order")
}



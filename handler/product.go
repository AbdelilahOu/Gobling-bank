package handler

import (
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Product struct {
	Repo *repository.ProductRepo
}

func (o *Product) Create(w http.ResponseWriter, r *http.Request) {

}

func (o *Product) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Products")
}

func (o *Product) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Product")
}

func (o *Product) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update Product")
}

func (o *Product) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete Product")
}

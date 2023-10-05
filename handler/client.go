package handler

import (
	"fmt"
	"net/http"

	"github.com/AbdelilahOu/GoThingy/repository"
)

type Client struct {
	Repo *repository.ClientRepo
}

func (o *Client) Create(w http.ResponseWriter, r *http.Request) {

}

func (o *Client) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Clients")
}

func (o *Client) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get Client")
}

func (o *Client) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update Client")
}

func (o *Client) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete Client")
}

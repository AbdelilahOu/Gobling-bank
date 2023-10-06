package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AbdelilahOu/GoThingy/model"
	"github.com/AbdelilahOu/GoThingy/repository"
	"github.com/google/uuid"
)

type Client struct {
	Repo *repository.ClientRepo
}

func (o *Client) Create(w http.ResponseWriter, r *http.Request) {
	// expected body struct
	var body struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
	}
	//error check (body doesnt match)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// create time
	now := time.Now().UTC()
	// create id
	id, err := uuid.NewUUID()
	// check error (creating uuid)
	if err != nil {
		fmt.Println("error generating the uuid", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create new client model
	client := model.Client{
		Id:         id,
		Firstname:  body.Firstname,
		Lastname:   body.Lastname,
		Phone:      body.Phone,
		Email:      body.Email,
		Created_at: &now,
	}
	// insert into db
	err = o.Repo.Insert(r.Context(), client)
	// check for error
	if err != nil {
		fmt.Println("error inserting the client", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// from model -> json
	res, err := json.Marshal(client)
	// check error
	if err != nil {
		fmt.Println("failed to marhsal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good return res
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Client) GetAll(w http.ResponseWriter, r *http.Request) {
	// pagination
	pageStr := r.URL.Query().Get("page")
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
	//
	const size = 20
	err = o.Repo.SelectAll(r.Context(), page, size)
	if err != nil {
		fmt.Println("error fetching client", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (o *Client) GetByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Client) UpdateByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Client) DeleteByID(w http.ResponseWriter, r *http.Request) {
}

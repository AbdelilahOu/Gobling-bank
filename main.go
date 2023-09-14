package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	// chi router
	router := chi.NewRouter()
	// use middleware
	router.Use(middleware.Logger)
	// get hello world
	router.Get("/hello-world",basicHandler)
	// create server
	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}
	// handle errors
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server : ", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello World"))
}
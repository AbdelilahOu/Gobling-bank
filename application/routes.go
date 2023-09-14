package application

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/AbdelilahOu/GoThingy/handler"
)

func loadRoutes() *chi.Mux {
	// create chi router
	router := chi.NewRouter()
	//  use logger middleware
	router.Use(middleware.Logger)
	// basinc route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// sub routes
	router.Route("/orders", loadOrderRoutes)
	// 
	return router
}

func loadOrderRoutes(router chi.Router) {
	// get order handler
	orderHandler := &handler.Order{}
	// attach routes
	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.GetAll)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)
}
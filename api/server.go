package api

import (
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/gin-gonic/gin"
)

// server serves http requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// create new server
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/entries", server.createEntry)
	router.GET("/entries/:id", server.getEntry)
	router.GET("/entries", server.listEntries)
	router.PUT("/entries/:id", server.updateEntry)
	router.DELETE("/entries/:id", server.deleteEntry)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// start server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

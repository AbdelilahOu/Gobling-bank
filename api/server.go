package api

import (
	"fmt"

	"github.com/AbdelilahOu/GoThingy/config"
	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// server serves http requests
type Server struct {
	config     config.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// create new server
func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(AuthMiddleware(tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.POST("/entries", server.createEntry)
	authRoutes.GET("/entries/:id", server.getEntry)
	authRoutes.GET("/entries", server.listEntries)
	authRoutes.PUT("/entries/:id", server.updateEntry)
	authRoutes.DELETE("/entries/:id", server.deleteEntry)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

// start server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

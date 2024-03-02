package api

import (
	"fmt"
	"net/http"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	fmt.Println("req:", req)
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// generate hash
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// create user
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, createUserResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	})
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get user
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// check password
	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	// generate token
	token, payload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, loginUserResponse{
		FullName: user.FullName,
		Username: payload.Username,
		Email:    user.Email,
		Token:    token,
	})
}

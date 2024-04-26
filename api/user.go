package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/AbdelilahOu/GoThingy/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
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
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// generate hash
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("generate hash password error")
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
				server.logger.Log.Error().Err(err).Msg("create user db error unique_violation")
				ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
				return
			}
		}
		server.logger.Log.Error().Err(err).Msg("create user error unique_violation")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// send verification email
	ops := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	err = server.taskDistributor.DistributTaskSendVerifyEmail(ctx, taskPayload, ops...)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("send verification email error")
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
	SessionID             uuid.UUID `json:"session_id"`
	Username              string    `json:"username"`
	FullName              string    `json:"full_name"`
	Email                 string    `json:"email"`
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get user
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Log.Error().Err(err).Msg("get user db error no row found")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get user error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// check password
	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("user login password check error")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	// generate token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("user login create access token error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// generate refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("user login create refresh token error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// create session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("user login create session error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, loginUserResponse{
		SessionID:             session.ID,
		FullName:              user.FullName,
		Username:              user.Username,
		Email:                 user.Email,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	})
}

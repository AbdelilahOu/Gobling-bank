package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/token"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/AbdelilahOu/GoThingy/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
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
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			FullName:       req.FullName,
			Email:          req.Email,
		},
		AfterCreate: func(user db.User) error {
			// send verification email
			ops := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			return server.taskDistributor.DistributTaskSendVerifyEmail(ctx, taskPayload, ops...)
		},
	}
	// create user
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			server.logger.Log.Error().Err(err).Msg("create user db error unique_violation")
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("create user error unique_violation")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, createUserResponse{
		Username: txResult.User.Username,
		Email:    txResult.User.Email,
		FullName: txResult.User.FullName,
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
		if errors.Is(err, db.ErrRecordNotFound) {
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
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.AccessTokenDuration)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("user login create access token error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// generate refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.RefreshTokenDuration)
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

type verifyEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

type verifyEmailRequest struct {
	EmailID    uuid.UUID `form:"id" binding:"required"`
	SecretCode string    `form:"secret_code" binding:"required"`
}

func (server *Server) verifyEmail(ctx *gin.Context) {
	var req verifyEmailRequest
	// validate the request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	//  verufy email
	result, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailID,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("couldnt verify email")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, verifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	})
}

type updateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	//
	value, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		server.logger.Log.Error().Msg("no payload object")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(fmt.Errorf("no payload object set")))
		return
	}
	payload := value.(*token.Payload)
	if !utils.HasPermission(payload.Role, []string{utils.BankerRole, utils.DepositorRole}) {
		server.logger.Log.Error().Msg("user role doesnt have permission")
		ctx.JSON(http.StatusConflict, utils.ErrorResponse(fmt.Errorf("user role doesnt have permission")))
		return
	}
	// get user
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get user for update db error no row found")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get user for update error no row found")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	//
	if payload.Role != utils.BankerRole && payload.Username != user.Username {
		server.logger.Log.Error().Msg("user name doesnt match payload")
		ctx.JSON(http.StatusConflict, utils.ErrorResponse(fmt.Errorf("username doesnt match")))
		return
	}
	// update user
	user, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
		Username: req.Username,
		FullName: pgtype.Text{
			Valid:  true,
			String: req.FullName,
		},
	})
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("update user db error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, user)
}

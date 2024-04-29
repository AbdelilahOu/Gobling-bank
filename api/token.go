package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// check if token is valid
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("verify token error")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	// get session
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get session db error")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get session error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	//
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	//
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("uncorrect session user")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	//
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	//
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	//
	// generate token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("create new access token error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	})
}

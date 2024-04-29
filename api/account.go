package api

import (
	"errors"
	"net/http"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/token"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get owner from token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// create account
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation || errCode == db.ForeignKeyViolation {
			server.logger.Log.Error().Err(err).Msg("create account db error foreign_key_violation or unique_violation")
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("create account error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,uuid"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	// validate the request
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get account
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get account db error no row found")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get account error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// get user from token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// check if user is the owner of the account
	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	// validate the request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get user from token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// get accounts
	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("list accounts error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	ID      uuid.UUID `uri:"id" binding:"required,uuid"`
	Balance int64     `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	// validate the request
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get account
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get account for update db error no row found")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get account for update error no row found")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// update account
	account, err = server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID:      account.ID,
		Balance: req.Balance,
	})
	if err != nil {
		server.logger.Log.Error().Err(err).Msg("update account db error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	// validate the request
	if err := ctx.ShouldBindUri(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get account
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get account for delete db error no row found")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("get account for delete error no row found")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// delete account
	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation || errCode == db.ForeignKeyViolation {
			server.logger.Log.Error().Err(err).Msg("delete account db error foreign_key_violation or unique_violation")
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("delete account error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, account)
}

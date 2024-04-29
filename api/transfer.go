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

type CreateTransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   uuid.UUID `json:"to_account_id"  binding:"required,uuid"`
	Amount        int64     `json:"amount"  binding:"required,gt=0"`
	Currency      string    `json:"currency"  binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Log.Error().Err(err).Msg("invalid request")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fromAccount, ok := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !ok {
		return
	}
	// get user from token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// check if user is the owner of the account
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("from account doesn't belong to authenticated user")
		server.logger.Log.Error().Err(err)
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	_, ok = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !ok {
		return
	}
	// run transaction
	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation || errCode == db.ForeignKeyViolation {
			server.logger.Log.Error().Err(err).Msg("create transfer db error foreign_key_violation or unique_violation")
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(err))
			return
		}
		server.logger.Log.Error().Err(err).Msg("create transfer error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return result
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID uuid.UUID, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			server.logger.Log.Error().Err(err).Msg("get account for transfer db error")
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return account, false
		}
		server.logger.Log.Error().Err(err).Msg("get account for transfer error")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := utils.ErrInvalidCurrency(account.ID, account.Currency, currency)
		server.logger.Log.Error().Err(err).Msg("account currency error")
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return account, false
	}

	return account, true
}

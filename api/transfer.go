package api

import (
	"database/sql"
	"net/http"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   uuid.UUID `json:"to_account_id"  binding:"required,uuid"`
	Amount        int64     `json:"amount"  binding:"required,gt=0"`
	Currency      string    `json:"currency"  binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	if !server.validAccountCurrency(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccountCurrency(ctx, req.ToAccountID, req.Currency) {
		return
	}
	// run transaction
	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return result
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccountCurrency(ctx *gin.Context, accountID uuid.UUID, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := utils.ErrInvalidCurrency(account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return false
	}

	return true
}

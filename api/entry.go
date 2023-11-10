package api

import (
	"net/http"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createEntryRequest struct {
	AccountID uuid.UUID `json:"account_id" binding:"required,uuid"`
	Amount    int64     `json:"amount" binding:"required"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
	// validate the request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// create entry
	entry, err := server.store.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, entry)
}

type getEntryRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,uuid"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest
	// validate the request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	// get entry
	entry, err := server.store.GetEntry(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	// return res
	ctx.JSON(http.StatusOK, entry)
}

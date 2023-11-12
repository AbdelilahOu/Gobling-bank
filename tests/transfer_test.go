package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/stretchr/testify/require"
)

func GenerateRandomTransfer(t *testing.T) db.Transfer {
	// generate accounts
	account1 := GenerateRandomAccount(t)
	account2 := GenerateRandomAccount(t)
	// transfer args
	arg := db.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomAmount(),
	}
	// create transfer
	transfer, err := testStore.CreateTransfer(context.Background(), arg)
	// check for errors
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	//
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	GenerateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	// generate transfer
	transfer := GenerateRandomTransfer(t)
	// get the generated transfer
	retrievedTransfer, err := testStore.GetTransfer(context.Background(), transfer.ID)
	// check for errors
	require.NoError(t, err)
	require.NotEmpty(t, retrievedTransfer)
	// check if feilds match
	require.Equal(t, transfer.ID, retrievedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, retrievedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, retrievedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, retrievedTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, retrievedTransfer.CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	// generate transfer
	transfer := GenerateRandomTransfer(t)
	// update transfer
	arg := db.UpdateTransferParams{
		ID:     transfer.ID,
		Amount: utils.RandomAmount(),
	}
	updatedTransfer, err := testStore.UpdateTransfer(context.Background(), arg)
	// check for errors
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)
	// check if feilds match
	require.Equal(t, transfer.ID, updatedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, updatedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, updatedTransfer.ToAccountID)
	require.Equal(t, arg.Amount, updatedTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, updatedTransfer.CreatedAt)

}

func TestDeleteTransfer(t *testing.T) {
	// generate transfer
	transfer := GenerateRandomTransfer(t)
	// delete transfer
	err := testStore.DeleteTransfer(context.Background(), transfer.ID)
	// check for errors
	require.NoError(t, err)
	// get transfer
	gotTransfer, err := testStore.GetTransfer(context.Background(), transfer.ID)
	// check for errors
	require.Error(t, err)
	require.Empty(t, gotTransfer)

}

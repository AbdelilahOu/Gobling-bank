package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/stretchr/testify/require"
)

func GenerateRandomEntry(t *testing.T) db.Entry {
	// get account
	account := GenerateRandomAccount(t)
	// get entry args
	arg := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomAmount(),
	}
	//
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	// check error
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	// check returned feilds
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	// check auto generated feilds
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry

}

func TestCreateEntry(t *testing.T) {
	GenerateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	// generate entry
	entry := GenerateRandomEntry(t)
	// get generated entry
	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	// check for errors
	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)
	// check returned feilds
	require.Equal(t, entry.ID, retrievedEntry.ID)
	require.Equal(t, entry.AccountID, retrievedEntry.AccountID)
	require.Equal(t, entry.Amount, retrievedEntry.Amount)
	require.Equal(t, entry.CreatedAt, retrievedEntry.CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry := GenerateRandomEntry(t)
	// update entry
	arg := db.UpdateEntryParams{
		ID:     entry.ID,
		Amount: utils.RandomAmount(),
	}
	//
	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	// check for errors
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)
	// check returned feilds
	require.Equal(t, arg.ID, updatedEntry.ID)
	require.Equal(t, arg.Amount, updatedEntry.Amount)
	require.Equal(t, entry.AccountID, updatedEntry.AccountID)
	require.Equal(t, entry.CreatedAt, updatedEntry.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	entry := GenerateRandomEntry(t)
	// delete entry
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	// check for errors
	require.NoError(t, err)
	// get deleted entry
	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	// check for errors
	require.Error(t, err)
	require.Empty(t, retrievedEntry)
}

package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	// create store
	store := db.NewStore(testDb)
	// create accounts
	account1 := GenerateRandomAccount(t)
	account2 := GenerateRandomAccount(t)
	// run a concurent transfer transactions
	n := 5
	amount := int64(10)
	// make channels to get data from go routines
	errs := make(chan error)
	results := make(chan db.TransferTxResult)
	// run transactions on a go routine
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	// checnking
	for i := 0; i < n; i++ {
		// check errors
		err := <-errs
		require.NoError(t, err)
		// check results
		result := <-results
		require.NotEmpty(t, result)
		// check for transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		// check for entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		// get entry from db
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		// get entry from db
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO : check accounts balance
	}
}

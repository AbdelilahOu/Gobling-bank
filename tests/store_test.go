package tests

import (
	"context"
	"fmt"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	// create accounts
	account1 := GenerateRandomAccount(t)
	account2 := GenerateRandomAccount(t)
	fmt.Println(">>>> before : ", account1.Balance, account2.Balance)
	// run a concurent transfer transactions
	n := 5
	amount := int64(1)
	// make channels to get data from go routines
	errs := make(chan error)
	results := make(chan db.TransferTxResult)
	// run transactions on a go routine
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}
	// checnking
	existed := make(map[int]bool)
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
		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		// get entry from db
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		// check accounts
		fromAccount := result.FromAccount
		fmt.Println("should not be empty:", fromAccount)
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">>>> on tx : ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1*amount, 2*amount, 3*amount, 4*amount, 5*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// check final ballence
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>>> after : ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
}

// func TestDeadLockTransfer(t *testing.T) {
// 	// create testStore
//
// 	// create accounts
// 	account1 := GenerateRandomAccount(t)
// 	account2 := GenerateRandomAccount(t)
// 	fmt.Println(">>>> before : ", account1.Balance, account2.Balance)
// 	// run a concurent transfer transactions
// 	n := 10
// 	amount := int64(10)
// 	// make channels to get data from go routines
// 	errs := make(chan error)
// 	// run transactions on a go routine
// 	for i := 0; i < n; i++ {
// 		FromAccountID := account1.ID
// 		ToAccountID := account2.ID
// 		if i%2 == 1 {
// 			FromAccountID = account2.ID
// 			ToAccountID = account1.ID
// 		}
// 		go func() {
// 			_, err := testStore.TransferTx(context.Background(), db.TransferTxParams{
// 				FromAccountID: FromAccountID,
// 				ToAccountID:   ToAccountID,
// 				Amount:        amount,
// 			})

// 			errs <- err
// 		}()
// 	}
// 	for i := 0; i < n; i++ {
// 		// check errors
// 		err := <-errs
// 		require.NoError(t, err)
// 		// get entry from db
// 	}
// 	// check final ballence
// 	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
// 	require.NoError(t, err)

// 	fmt.Println(">>>> after : ", updatedAccount1.Balance, updatedAccount2.Balance)

// 	require.Equal(t, updatedAccount1.Balance, account1.Balance)
// 	require.Equal(t, updatedAccount2.Balance, account2.Balance)
// }

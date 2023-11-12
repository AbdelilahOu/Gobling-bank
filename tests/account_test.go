package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/stretchr/testify/require"
)

func GenerateRandomAccount(t *testing.T) db.Account {
	user := GenerateRandomUser(t)

	arg := db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomAmount(),
		Currency: utils.RandomCurrency(),
	}
	// create account
	account, err := testStore.CreateAccount(context.Background(), arg)
	// check error
	require.NoError(t, err)
	require.NotEmpty(t, account)
	// check returned feilds
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	// check auto generated feilds
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	GenerateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// generate account
	account := GenerateRandomAccount(t)
	// retrieved account
	retrievedAccount, err := testStore.GetAccount(context.Background(), account.ID)
	// check errors
	require.NoError(t, err)
	require.NotEmpty(t, retrievedAccount)
	// check returned feilds
	require.Equal(t, account.ID, retrievedAccount.ID)
	require.Equal(t, account.Owner, retrievedAccount.Owner)
	require.Equal(t, account.Balance, retrievedAccount.Balance)
	require.Equal(t, account.Currency, retrievedAccount.Currency)
	require.Equal(t, account.CreatedAt, retrievedAccount.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account := GenerateRandomAccount(t)
	arg := db.UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomAmount(),
	}
	updatedAcc, err := testStore.UpdateAccount(context.Background(), arg)
	// check errors
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	// check returned feilds
	require.Equal(t, account.ID, updatedAcc.ID)
	require.Equal(t, account.Owner, updatedAcc.Owner)
	require.Equal(t, arg.Balance, updatedAcc.Balance)
	require.Equal(t, account.Currency, updatedAcc.Currency)
	require.Equal(t, account.CreatedAt, updatedAcc.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := GenerateRandomAccount(t)
	err := testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

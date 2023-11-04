package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/stretchr/testify/require"
)

func GenerateRandomAccount(t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomAmount(),
		Currency: utils.RandomCurrency(),
	}
	// create account
	account, err := testQueries.CreateAccount(context.Background(), arg)
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
	account1 := GenerateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// check errors
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// check returned feilds
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

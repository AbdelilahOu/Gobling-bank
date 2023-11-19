package tests

import (
	"context"
	"testing"

	db "github.com/AbdelilahOu/GoThingy/db/sqlc"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/stretchr/testify/require"
)

func GenerateRandomUser(t *testing.T) db.User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	arg := db.CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	// create user
	user, err := testStore.CreateUser(context.Background(), arg)
	// check error
	require.NoError(t, err)
	require.NotEmpty(t, user)
	// check returned feilds
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	// check auto generated feilds
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	GenerateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// generate user
	user := GenerateRandomUser(t)
	// retrieved user
	retrievedUser, err := testStore.GetUser(context.Background(), user.Username)
	// check errors
	require.NoError(t, err)
	require.NotEmpty(t, retrievedUser)
	// check returned feilds
	require.Equal(t, user.Username, retrievedUser.Username)
	require.Equal(t, user.HashedPassword, retrievedUser.HashedPassword)
	require.Equal(t, user.FullName, retrievedUser.FullName)
	require.Equal(t, user.Email, retrievedUser.Email)
}

// func TestUpdateUser(t *testing.T) {
// 	user := GenerateRandomUser(t)
// 	arg := db.UpdateUserParams{
// 		ID:      user.ID,
// 		Balance: utils.RandomAmount(),
// 	}
// 	updatedAcc, err := testStore.UpdateUser(context.Background(), arg)
// 	// check errors
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedAcc)
// 	// check returned feilds
// 	require.Equal(t, user.ID, updatedAcc.ID)
// 	require.Equal(t, user.Username, updatedAcc.Username)
// 	require.Equal(t, arg.Balance, updatedAcc.Balance)
// 	require.Equal(t, user.Currency, updatedAcc.Currency)
// 	require.Equal(t, user.CreatedAt, updatedAcc.CreatedAt)
// }

// func TestDeleteUser(t *testing.T) {
// 	user := GenerateRandomUser(t)
// 	err := testStore.DeleteUser(context.Background(), user.ID)
// 	require.NoError(t, err)
// 	user2, err := testStore.GetUser(context.Background(), user.ID)
// 	require.Error(t, err)
// 	require.Empty(t, user2)
// }

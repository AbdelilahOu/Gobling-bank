package tests

import (
	"context"
	"database/sql"
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

func TestUpdateUserFullname(t *testing.T) {
	oldUser := GenerateRandomUser(t)
	fullName := utils.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			Valid:  true,
			String: fullName,
		},
	})
	// check errors
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	// check returned feilds
	require.Equal(t, fullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := GenerateRandomUser(t)
	email := utils.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			Valid:  true,
			String: email,
		},
	})
	// check errors
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	// check returned feilds
	require.Equal(t, email, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := GenerateRandomUser(t)
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			Valid:  true,
			String: hashedPassword,
		},
	})
	// check errors
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	// check returned feilds
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}

func TestUpdateUser(t *testing.T) {
	oldUser := GenerateRandomUser(t)
	fullName := utils.RandomOwner()
	email := utils.RandomEmail()
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			Valid:  true,
			String: hashedPassword,
		},
		FullName: sql.NullString{
			Valid:  true,
			String: fullName,
		},
		Email: sql.NullString{
			Valid:  true,
			String: email,
		},
	})
	// check errors
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	// check returned feilds
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, email, updatedUser.Email)
	require.Equal(t, fullName, updatedUser.FullName)
}

// func TestDeleteUser(t *testing.T) {
// 	user := GenerateRandomUser(t)
// 	err := testStore.DeleteUser(context.Background(), user.ID)
// 	require.NoError(t, err)
// 	user2, err := testStore.GetUser(context.Background(), user.ID)
// 	require.Error(t, err)
// 	require.Empty(t, user2)
// }

package tests

import (
	"testing"
	"time"

	"github.com/AbdelilahOu/GoThingy/token"
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := -time.Minute

	createdToken, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, createdToken)

	payload, err := maker.VerifyToken(createdToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlg(t *testing.T) {
	maker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err := token.NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	signedToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(signedToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrInvalidToken.Error())
	require.Nil(t, payload)
}

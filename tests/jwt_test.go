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

	token, _, err := maker.CreateToken(username, utils.DepositorRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	verifiedPayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, verifiedPayload)

	require.NotZero(t, verifiedPayload.ID)
	require.Equal(t, username, verifiedPayload.Username)
	require.Equal(t, utils.DepositorRole, verifiedPayload.Role)
	require.WithinDuration(t, issuedAt, verifiedPayload.IssuedAt, time.Second)
	require.WithinDuration(t, expredAt, verifiedPayload.ExpiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := -time.Minute

	createdToken, _, err := maker.CreateToken(username, utils.DepositorRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, createdToken)

	verifiedPayload, err := maker.VerifyToken(createdToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrExpiredToken.Error())
	require.Nil(t, verifiedPayload)
}

func TestInvalidJWTTokenAlg(t *testing.T) {
	maker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err := token.NewPayload(utils.RandomOwner(), utils.DepositorRole, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	signedToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(signedToken)
	require.Error(t, err)
	require.EqualError(t, err, token.ErrInvalidToken.Error())
	require.Nil(t, payload)
}

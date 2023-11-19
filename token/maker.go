package token

import "time"

type Maker interface {
	// crreate a sign
	CreateToken(username string, duration time.Duration) (string, error)
	// verify sign
	VerifyToken(token string) (*Payload, error)
}

package token

import "github.com/o1egl/paseto"

type PastoMaker struct {
	pasto        *paseto.V2
	symmetricKey []byte
}

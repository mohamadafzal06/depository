package service

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Number int64
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

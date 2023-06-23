package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mohamadafzal06/depository/param"
)

type AuthConfig struct {
	SignKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

type Auth struct {
	config AuthConfig
}

func NewAuth(cfg AuthConfig) Auth {
	return Auth{
		config: cfg,
	}
}

func (a Auth) CreateAccessToken(req param.LoginRequest) (param.LoginResponse, error) {
	tokenString, err := a.createToken(req.Number, a.config.AccessSubject, a.config.AccessExpirationTime)
	if err != nil {
		return param.LoginResponse{TokenString: "", Status: param.LoginUnsuccessful}, fmt.Errorf("cannot login: %w", err)
	}

	return param.LoginResponse{TokenString: tokenString, Status: param.LoginSuccessful}, nil
}

func (a Auth) CreateRefreshToken(req param.LoginRequest) (param.LoginResponse, error) {

	tokenString, err := a.createToken(req.Number, a.config.RefreshSubject, a.config.RefreshExpirationTime)
	if err != nil {
		return param.LoginResponse{TokenString: "", Status: param.LoginUnsuccessful}, fmt.Errorf("cannot login: %w", err)
	}

	return param.LoginResponse{TokenString: tokenString, Status: param.LoginSuccessful}, nil
}

func (a Auth) ParseToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (a Auth) createToken(number int64, subject string, expireDuration time.Duration) (string, error) {

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		Number: number,
	}

	// TODO - add sign method to config
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(a.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

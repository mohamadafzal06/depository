package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mohamadafzal06/depository/param"
	"github.com/mohamadafzal06/depository/service"
)

func permissioinDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, HandlerErr{Error: "permissioin denied"})
}

func JWTMiddleware(hrFunc http.HandlerFunc, srv *service.Depository, authSrv *service.Auth, authCfg *service.AuthConfig) http.HandlerFunc {
	fmt.Println("calling JWT auth middleware")

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		parts := strings.Split(authHeader, " ")
		tokenString := parts[1]

		// token validation
		token, err := validateJWT(tokenString, authCfg)
		if err != nil {
			permissioinDenied(w)
		}

		if !token.Valid {
			permissioinDenied(w)
		}

		var req param.GetAccountByNumberRequest
		number := getNumber(r)
		if number == -1 {
			permissioinDenied(w)
			return
		}

		req.Number = number

		account, err := srv.GetAccountByNumber(r.Context(), req)
		if err != nil {
			permissioinDenied(w)
			return
		}

		// parsing token for getting account number
		claims, _ := authSrv.ParseToken(tokenString)
		if account.Number != claims.Number {
			permissioinDenied(w)
			return
		}

		hrFunc(w, r)
	}

}

func validateJWT(tokenString string, cfg *service.AuthConfig) (*jwt.Token, error) {
	secret := cfg.SignKey

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

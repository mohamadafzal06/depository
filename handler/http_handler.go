package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/mohamadafzal06/depository/param"
	"github.com/mohamadafzal06/depository/service"
)

type HandlerErr struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v ...interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, HandlerErr{Error: err.Error()})
		}
	}
}

type Handler struct {
	listenAddr string
	service    *service.Depository
}

func New(lAddr string, srv *service.Depository) *Handler {
	return &Handler{
		listenAddr: lAddr,
		service:    srv,
	}
}

func (h *Handler) Run() {
	router := http.NewServeMux()

	router.HandleFunc("/account", makeHTTPHandleFunc(h.handleAccount))
	router.HandleFunc("/account/{number}", JWTMiddleware(makeHTTPHandleFunc(h.handleGetAccount), h.service))
	router.HandleFunc("/account/remove/{number}", JWTMiddleware(makeHTTPHandleFunc(h.handleDeleteAccount), h.service))
	router.HandleFunc("/transfer", JWTMiddleware(makeHTTPHandleFunc(h.handleTransfer), h.service))

	log.Printf("Handler is running on port: %s\n", h.listenAddr)

	http.ListenAndServe(h.listenAddr, router)
}

func (s *Handler) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetAccount(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (h *Handler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid method")
	}

	var createdAccountReq param.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&createdAccountReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("account creation failed."))
		return fmt.Errorf("cannot bind requst body to request param: %w", err)
	}

	createdAccountResponse, err := h.service.CreateAccount(r.Context(), createdAccountReq)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, []byte("account creation failed."))
	}

	tokenString, err := createJWT(&createdAccountResponse)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, []byte("authentication failed."))
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	return WriteJSON(w, http.StatusOK, createdAccountResponse)
}

func (h *Handler) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method")
	}

	var req param.GetAccountByNumberRequest
	number := getNumber(r)
	if number != -1 {
		req.Number = getNumber(r)
	} else {
		return WriteJSON(w, http.StatusBadRequest, []byte("the number is not valid"))
	}

	response, err := h.service.GetAccountByNumber(r.Context(), req)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, []byte("cannot get account with this number."))
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		return fmt.Errorf("invalid method")
	}

	var req param.DeleteAccountRequest
	number := getNumber(r)

	if number != -1 {
		req.Number = number
	} else {
		return WriteJSON(w, http.StatusBadRequest, []byte("the number is not valid"))
	}

	err := h.service.DeleteAccount(r.Context(), req)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, []byte("cannot delete account with this number."))
	}

	return WriteJSON(w, http.StatusOK, []byte("the account has been removed successully."))
}

func (h *Handler) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid method")
	}
	var req param.TransferAmountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("cannot bind the request body: %w", err)
	}

	defer r.Body.Close()

	response, err := h.service.TransferAmount(r.Context(), req)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, response)
	}

	return WriteJSON(w, http.StatusOK, response)
}

func permissioinDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, HandlerErr{Error: "permissioin denied"})
}

func JWTMiddleware(hrFunc http.HandlerFunc, srv *service.Depository) http.HandlerFunc {
	fmt.Println("calling JWT auth middleware")

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
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

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["number"].(float64)) {
			permissioinDenied(w)
			return
		}

		hrFunc(w, r)
	}

}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func createJWT(createdAccount *param.CreateAccountResponse) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["number"] = createdAccount.Number

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func getNumber(r *http.Request) int64 {
	vars := mux.Vars(r)
	nString := vars["number"]
	n, err := strconv.Atoi(nString)
	if err != nil {
		return -1
	}
	return int64(n)
}

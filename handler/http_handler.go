package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/account/{number}", makeHTTPHandleFunc(h.handleGetAccount))
	router.HandleFunc("/account/remove/{number}", makeHTTPHandleFunc(h.handleDeleteAccount))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(h.handleTransfer))

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

	return WriteJSON(w, http.StatusOK, createdAccountResponse)
}

func (h *Handler) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method")
	}
	vars := mux.Vars(r)
	nString := vars["number"]
	var req param.GetAccountByNumberRequest
	n, err := strconv.Atoi(nString)
	if err != nil {
		return fmt.Errorf("invalid number given: %w", err)
	}
	req.Number = int64(n)

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

	vars := mux.Vars(r)
	nString := vars["number"]
	var req param.DeleteAccountRequest
	n, err := strconv.Atoi(nString)
	if err != nil {
		return fmt.Errorf("invalid number given: %w", err)
	}
	req.Number = int64(n)

	err = h.service.DeleteAccount(r.Context(), req)
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ServerErr struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ServerErr{Error: err.Error()})
		}
	}
}

type Server struct {
	listenAddr string
}

func New(lAddr string) *Server {
	return &Server{
		listenAddr: lAddr,
	}
}

func (s *Server) Run() {
	router := http.NewServeMux()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	//	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Printf("Server is running on port: %s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *Server) handleAccount(w http.ResponseWriter, r *http.Request) error {
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

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("mamreza", "afzal")

	return WriteJSON(w, http.StatusOK, account)
}

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

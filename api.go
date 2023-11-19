package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))

	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleAccountById), s.store))

	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.handleTransfer))

	log.Println("Starting APIServer...")

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		loginRequest := &LoginRequest{}

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			return err
		}

		acc, err := s.store.GetAccountByEmail(loginRequest.Email)

		if err != nil {
			return err
		}

		fmt.Println("Account =>", acc)

		// compare password
		if err := bcrypt.CompareHashAndPassword([]byte(acc.HashPassword), []byte(loginRequest.Password)); err != nil {
			return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid credentials"})
		}

		// create JWT token
		
	tokenString, err := createJWTToken(acc)

	if err != nil {
		return err
	}

	fmt.Println("JWT Token string =>", tokenString)

	payload := &LoginResponse{
		Token: tokenString,
		Email: acc.Email,
		Id: acc.ID,
	}
		return WriteJSON(w, http.StatusOK, payload)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, e *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)

		if err != nil {
			return err
		}

		account, err := s.store.GetAccountById(id)

		if err != nil {
			return WriteJSON(w, http.StatusNotFound, ApiError{Error: "User not found"})
		}

		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	// OR
	// createAccountReq := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := NewAccount(
		req.FirstName, 
		req.LastName, 
		req.Email, 
		req.Password,
	)

	
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, Iderr := getID(r)

	if Iderr != nil {
		return Iderr
	}

	err := s.store.DeleteAccount(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := &TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferRequest)
}

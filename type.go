package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	HashPassword string    `json:"hashPassword"`
	Number       int64     `json:"number"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Id 		int`json:"id"`
}

func NewAccount(firstName, lastName, email, password string) (*Account, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		Number:       int64(rand.Intn(100000001)),
		Email:        email,
		HashPassword: string(hashPassword),
		CreatedAt:    time.Now().UTC(),
	}, nil
}

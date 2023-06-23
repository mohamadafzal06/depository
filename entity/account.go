package entity

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                uint64    `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

// TODO: can be replaced with uuid for Account's Number
func RandomNumber() int64 {
	rand.Seed(time.Now().UnixNano())

	// generate an 8-digit random number between 10000000 and 99999999
	id := rand.Int63n(90000000) + 10000000
	return id
}

func NewAccount(fn, ln, password string, balance int64) (*Account, error) {
	encpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot create new account: %w", err)
	}
	n := RandomNumber()
	id := uint64(rand.Intn(10000))
	return &Account{
		ID:                id,
		FirstName:         fn,
		LastName:          ln,
		EncryptedPassword: string(encpass),
		Number:            n,
		Balance:           balance,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

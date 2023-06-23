package entity

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewAccount(t *testing.T) {
	fn := "John"
	ln := "Doe"
	password := "mypassword"

	account, err := NewAccount(fn, ln, password)

	if err != nil {
		t.Errorf("unexpected error while creating new account: %s", err.Error())
	}

	// Check if values are set correctly
	if account.FirstName != fn {
		t.Errorf("expected first name to be %s, but got %s", fn, account.FirstName)
	}
	if account.LastName != ln {
		t.Errorf("expected last name to be %s, but got %s", ln, account.LastName)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(password)); err != nil {
		t.Errorf("expected password to be encrypted, but got %s", account.EncryptedPassword)
	}
	if account.Number < 10000000 || account.Number > 99999999 {
		t.Errorf("expected account number to be between 10000000 and 99999999, but got %d", account.Number)
	}
}

package repository

import "github.com/mohamadafzal06/depository/entity"

type Repository interface {
	CreateAccount(acc *entity.Account) error
	DeleteAccount(id int) error
	UpdateAccount(acc *entity.Account) error
	GetAccountByID(id int) (*entity.Account, error)
}

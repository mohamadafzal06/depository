package repository

import (
	"context"

	"github.com/mohamadafzal06/depository/entity"
)

type Repository interface {
	CreateAccount(ctx context.Context, acc *entity.Account) (int64, error)
	DeleteAccount(ctx context.Context, number int64) error
	TransferAmount(ctx context.Context, from, to, amount int64) error
	GetAccountByNumber(ctx context.Context, number int64) (*entity.Account, error)
}

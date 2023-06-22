package service

import (
	"fmt"

	"context"

	"github.com/mohamadafzal06/depository/entity"
	"github.com/mohamadafzal06/depository/param"
	"github.com/mohamadafzal06/depository/repository"
)

type Depository struct {
	repo repository.Repository
}

func NewDepository(r repository.Repository) *Depository {
	return &Depository{
		repo: r,
	}
}

func (s *Depository) CreateAccount(ctx context.Context, req param.CreateAccountRequest) (param.CreateAccountResponse, error) {
	acc := entity.Account{
		FirstName: req.FistName,
		LastName:  req.LastName,
		Balance:   req.Balance,
	}

	number, err := s.repo.CreateAccount(ctx, &acc)
	if err != nil {
		return param.CreateAccountResponse{}, fmt.Errorf("cannot create this account: %w", err)
	}

	response := param.CreateAccountResponse{
		FistName: req.FistName,
		LastName: req.LastName,
		Number:   number,
	}

	return response, nil
}

func (s *Depository) GetAccountByNumber(ctx context.Context, req param.GetAccountByNumberRequest) (param.GetAccountByNumberResponse, error) {
	acc, err := s.repo.GetAccountByNumber(ctx, req.Number)
	if err != nil {
		return param.GetAccountByNumberResponse{}, fmt.Errorf("cannot get account by this number: %w", err)
	}
	response := param.GetAccountByNumberResponse{
		FistName:  acc.FirstName,
		LastName:  acc.LastName,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
	}

	return response, nil
}

func (s *Depository) DeleteAccount(ctx context.Context, req param.DeleteAccountRequest) error {
	err := s.repo.DeleteAccount(ctx, req.Number)
	if err != nil {
		return fmt.Errorf("cannot get account by this number: %w", err)
	}

	return nil
}

func (s *Depository) TransferAmount(ctx context.Context, req param.TransferAmountRequest) (param.TransferAmountResponse, error) {
	err := s.repo.TransferAmount(ctx, req.FromAccount, req.ToAccount, req.Amount)
	if err != nil {
		return param.TransferAmountResponse{Status: param.Unsuccessful}, fmt.Errorf("transfer money failed: %w", err)
	}
	return param.TransferAmountResponse{Status: param.Successful}, nil
}

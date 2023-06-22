package param

import "time"

type TransferStatus string

const (
	Successful   TransferStatus = "Successful"
	Unsuccessful                = "Unsuccessful"
)

type CreateAccountRequest struct {
	FistName string `json:"fist_name"`
	LastName string `json:"last_name"`
	Balance  int64  `json:"balance"`
}
type CreateAccountResponse struct {
	FistName string `json:"fist_name"`
	LastName string `json:"last_name"`
	Number   int64  `json:"number"`
}

type GetAccountByNumberRequest struct {
	Number int64 `json:"number"`
}
type GetAccountByNumberResponse struct {
	FistName  string    `json:"fist_name"`
	LastName  string    `json:"last_name"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteAccountRequest struct {
	Number int64 `json:"number"`
}

type TransferAmountRequest struct {
	FromAccount int64 `json:"from_account"`
	ToAccount   int64 `json:"to_account"`
	Amount      int64 `json:"amount"`
}
type TransferAmountResponse struct {
	Status TransferStatus `json:"status"`
}

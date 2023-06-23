package param

import "time"

type TransferStatus string

const (
	Successful   TransferStatus = "Successful"
	Unsuccessful TransferStatus = "Unsuccessful"
)

type LoginStatus string

const (
	LoginSuccessful   LoginStatus = "Login Successful"
	LoginUnsuccessful LoginStatus = "Login Unsuccessful"
)

type CreateAccountRequest struct {
	FistName string `json:"fist_name"`
	LastName string `json:"last_name"`
	Password string `json:"password"`
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

type CreateTokenRequst struct {
	Number int64
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	TokenString string
	Status      LoginStatus
}

type PassCheckRespone struct {
	Truly bool
}

package account

import "github.com/pkg/errors"

var ErrAuthorizationRequired = errors.New("Authorization required")
var ErrInvalidAmount = errors.New("Invalid input amount")
var ErrOverdrawn = errors.New("Overdrawn")

/*
Delegate is a delegate of account
ID return the id of the user
Withdraw do withdraw, return fee, balance after withdraw
Deposit do deposit, return balance after deposit
History return history of transactions(withdraw and deposit)
Balance return current balance in the account
 */
type Delegate interface {
	ID() string
	Withdraw(amount float64) (float64, float64, error)
	Deposit(amount float64) (float64, error)
	History() ([]Transaction, error)
	Balance() (float64, error)
}

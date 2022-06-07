package cashier

import "github.com/pkg/errors"

var ErrNotEnoughCash = errors.New("Not Enough Cash")
var ErrInvalidRequest = errors.New("Invalid Request")
var ErrUnknown = errors.New("Unknown error")

/*
Cashier manage the money in the machine
WithdrawValidateAndAdjust check if the amount for withdraw is valid and return the real number can withdraw
Withdraw do withdraw operation, must call WithdrawValidateAndAdjust before this function, or may return error
Deposit do deposit operation
Total return amount of cash in the machine
 */
type Cashier interface {
	WithdrawValidateAndAdjust(amount float64) (float64, error)
	Withdraw(amount float64) error
	Deposit(amount float64) error
	Total() float64
}

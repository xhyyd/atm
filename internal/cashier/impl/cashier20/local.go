package cashier20

import (
	"github.com/pkg/errors"
	"github.com/xhyyd/atm/internal/account"
	"github.com/xhyyd/atm/internal/cashier"
	"math"
)

type impl struct {
	amount float64
}

var errInvalidFormat = errors.New("Invalid input, amount must be positive and multiple of 20.")

func (i impl) WithdrawValidateAndAdjust(amount float64) (float64, error) {
	if amount <= 0 || amount != math.Trunc(amount) || int(amount) % 20 != 0 {
		return 0, errInvalidFormat
	}

	if amount > i.amount {
		amount = i.amount
	}

	if amount != math.Trunc(amount) || int(amount) % 20 != 0 {
		amount = float64(int(amount) / 20 * 20)
	}
	return amount, nil
}

func (i *impl) Withdraw(amount float64) error {
	if amount > i.amount {
		return cashier.ErrNotEnoughCash
	}
	if amount != math.Trunc(amount) || int(amount) % 20 != 0 {
		return cashier.ErrInvalidRequest
	}
	i.amount -= amount
	return nil
}

func (i *impl) Deposit(amount float64) error {
	if amount <= 0 {
		return account.ErrInvalidAmount
	}
	i.amount += amount
	return nil
}

func (i impl) Total() float64 {
	return math.Round(i.amount * 100) / 100
}

func NewCashierWithMoney(val float64) *impl {
	return &impl{
		amount: val,
	}
}
package local

import (
	"github.com/xhyyd/atm/internal/account"
	"math"
	"time"
)



type delegate struct {
	item *accountItem
	loginTime time.Time
	lastTrans string
}

var a account.Delegate = (*delegate)(nil)

func (d delegate) ID() string {
	return d.item.account.ID
}

func (d *delegate) Withdraw(amount float64) (float64, float64, error) {
	if d == nil {
		return 0, d.balance(), account.ErrAuthorizationRequired
	}

	if amount <= 0 {
		return 0, d.balance(), account.ErrInvalidAmount
	}

	var fee float64
	if d.item.account.Balance < 0 {
		return 0, d.balance(), account.ErrOverdrawn
	}

	if d.item.account.Balance < amount {
		fee = 5.0
	}

	d.applyTransaction(account.Transaction{
		Time: time.Now(),
		Value: -amount - fee,
	})
	return fee, d.balance(), nil
}

func (d *delegate) Deposit(amount float64) (float64, error) {
	if d == nil {
		return 0, account.ErrAuthorizationRequired
	}
	d.applyTransaction(account.Transaction{
		Time: time.Now(),
		Value: amount,
	})
	return d.balance(), nil
}

func (d *delegate) History() ([]account.Transaction, error) {
	if d == nil {
		return nil, account.ErrAuthorizationRequired
	}
	return d.item.history, nil
}

func (d *delegate) balance() float64 {
	return math.Round(d.item.account.Balance * 100) / 100
}

func (d *delegate) Balance() (float64, error) {
	if d == nil {
		return 0, account.ErrAuthorizationRequired
	}
	return d.balance(), nil
}

func (d *delegate) applyTransaction(trans account.Transaction) {
	d.item.account.Balance += trans.Value
	trans.Balance = d.item.account.Balance
	d.item.history = append([]account.Transaction{trans}, d.item.history...)
}

func newDelegate(item *accountItem) *delegate {
	return &delegate{
		item:item,
		loginTime: time.Now(),
	}
}

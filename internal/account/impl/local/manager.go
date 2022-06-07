package local

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xhyyd/atm/internal/account"
	"strconv"
)

type accountItem struct {
	account account.Account
	history []account.Transaction
}

func newAccount(id string, pin string, balance float64) *accountItem {
	return &accountItem{
		account:account.Account {
			ID: id,
			Pin: pin,
			Balance: balance,
		},
		history: nil,
	}
}

type manager struct {
	database map[string]*accountItem
}



func (m manager) Authorize(id, pin string) (account.Delegate, error) {
	d, ok := m.database[id]
	if !ok {
		return nil, account.ErrNotFound
	}
	if pin != d.account.Pin {
		return nil, account.ErrInvalidPin
	}

	return newDelegate(d), nil
}

// NewManagerWithData create manager and init with data
// id, pin, balance
func NewManagerWithData(data [][]string) (account.Manager, error) {
	database := map[string]*accountItem{}

	for i, d := range data {
		name := d[0]
		pin := d[1]
		balance, err := strconv.ParseFloat(d[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse %d line", i))
		}
		database[name] = newAccount(name, pin, balance)
	}

	return &manager{
		database: database,
	}, nil
}


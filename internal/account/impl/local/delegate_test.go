package local

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xhyyd/atm/internal/account"
	"testing"
)

func TestDelegate_Balance(t *testing.T) {
	type Testcase struct {
		balance float64
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{10.24}, {90000.55}, { 0},
	}
	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			dl := newDelegate(&accountItem{account: account.Account{Balance: v.balance}})
			b, err := dl.Balance()
			assert.Nil(t, err)
			assert.Equal(t, v.balance, b)
		})
	}
}

func TestDelegate_Deposit(t *testing.T) {
	type Testcase struct {
		balance float64
		deposit float64
		balanceAfterDeposit float64
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{10.24, 0, 10.24},
		{0, 90000.55, 90000.55},
		{0, 0, 0},
		{60.00, 55.24, 115.24},
	}
	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			dl := newDelegate(&accountItem{account: account.Account{Balance: v.balance}})
			dl.Deposit(v.deposit)
			b, err := dl.Balance()
			assert.Nil(t, err)
			assert.Equal(t, v.balanceAfterDeposit, b)
		})
	}
}

func TestDelegate_Withdraw(t *testing.T) {
	type Testcase struct {
		balance float64
		withdraw float64
		balanceAfterWithdraw float64
		fee float64
		err error
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{10.24, 0, 10.24, 0, account.ErrInvalidAmount},
		{0, 90000.55, -90005.55, 5, nil},
		{0, 0, 0, 0, account.ErrInvalidAmount},
		{60.00, 55.24, 4.76, 0,nil},
		{60.00, 100, -45, 5,nil},
	}
	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			dl := newDelegate(&accountItem{account: account.Account{Balance: v.balance}})
			fee, balance, err := dl.Withdraw(v.withdraw)
			assert.Equal(t, v.err, err)
			assert.Equal(t, v.fee, fee)
			assert.Equal(t, v.balanceAfterWithdraw, balance)
		})
	}
}

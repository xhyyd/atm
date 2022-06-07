package atm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xhyyd/atm/internal/account/impl/local"
	"github.com/xhyyd/atm/internal/cashier/impl/cashier20"
	"github.com/xhyyd/atm/internal/ui/impl/mock"
	"os"
	"strings"
	"testing"
	"time"
)

func defaultATM() (*ATM, *mock.MockUI) {
	cashier := cashier20.NewCashierWithMoney(10000)
	am, err := local.NewManagerWithData([][]string{
		{"2859459814","7386","10.24"},
		{"1434597300","4557","90000.55"},
		{"7089382418","0075","0.00"},
		{"2001377812","5950","60.00"},
	})
	if err != nil {
		fmt.Println("error in initialize account manager:", err.Error())
		os.Exit(1)
	}
	ui := mock.NewUI()
	return NewATM(am, cashier, ui), ui
}

func TestATM_Authorize(t *testing.T) {
	type Testcase struct {
		id string
		pin string
		err error
		output []string
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{"2859459814","7386", nil, []string{"2859459814 successfully authorized."}},
		{"7089382418","7381", nil, []string{"Authorization failed."}},
		{"7089382413","5950", nil, []string{"Authorization failed."}},
	}

	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			m, ui := defaultATM()
			err := m.Authorize(v.id, v.pin)
			assert.Equal(t, v.err, err)
			checkStringArray(v.output, ui.Lines, t)
		})
	}
}

func TestATM_Balance(t *testing.T) {
	type Testcase struct {
		id string
		pin string
		output []string
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{"2859459814","7386", []string{"Current balance: 10.24"}},
		{"1434597300","4557", []string{"Current balance: 90000.55"}},
		{"7089382418","0075",[]string{"Current balance: 0.00"}},
		{"2001377812","5950", []string{"Current balance: 60.00"}},
	}

	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			m, ui := defaultATM()
			err := m.Authorize(v.id, v.pin)
			assert.Nil(t, err)
			ui.Clear()
			err = m.Balance()
			assert.Nil(t, err)
			checkStringArray(v.output, ui.Lines, t)
		})
	}
}

func TestATM_Wait(t *testing.T) {
	m, ui := defaultATM()
	err := m.Authorize("1434597300", "4557")
	assert.Nil(t, err)
	checkStringArray([]string{"1434597300 successfully authorized."}, ui.Lines, t)
	ui.Clear()
	time.Sleep(2 * time.Minute + 1 * time.Second)
	err = m.Deposit(20)
	assert.Nil(t, err)
	checkStringArray([]string{"Authorization required."}, ui.Lines, t)
}

func TestATM_WithdrawDepositHistory(t *testing.T) {
	type Testcase struct {
		operation string
		val interface{}
		output []string
	}
	// todo add more testcases
	type Testsuite []Testcase
	ts := Testsuite {
		{"authorize", []string{"1434597300", "4557"}, []string{"1434597300 successfully authorized."}},
		{"withdraw", 0.0,[]string{"Invalid input, amount must be positive and multiple of 20."}},
		{"deposit", 0.0,[]string{"Invalid input, amount must be positive."}},
		{"history", nil,[]string{"No history found"}},
		{"withdraw", 20.0,[]string{"Amount dispensed: $20.00", "Current balance: 89980.55"}},
		{"withdraw", 90000.0,[]string{"Unable to dispense full amount requested at this time.", "Amount dispensed: $9980.00", "Current balance: 80000.55"}},
		{"withdraw", 20.0,[]string{"Unable to process your withdrawal at this time."}},
		{"deposit", 100.0,[]string{"Current balance: 80100.55"}},
		{"history", nil,[]string{"100.00 80100.55", "-9980.00 80000.55", "-20.00 89980.55"}},
		{"authorize", []string{"2001377812", "5950"}, []string{"2001377812 successfully authorized."}},
		{"withdraw", 80.0,[]string{"Amount dispensed: $80.00", "You have been charged an overdraft fee of $5.00. Current balance: -25.00"}},
		{"withdraw", 20.0,[]string{"Your account is overdrawn. Current balance: -25.00"}},
	}
	m, ui := defaultATM()
	err := m.Authorize("1434597300","4557")
	ui.Clear()
	assert.Nil(t, err)
	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("Test_%d" , i), func (t *testing.T){
			ui.Clear()
			if v.operation == "authorize" {
				err = m.Authorize(v.val.([]string)[0], v.val.([]string)[1])
			}else if v.operation == "withdraw" {
				err = m.Withdraw(v.val.(float64))
			} else if v.operation == "deposit" {
				err = m.Deposit(v.val.(float64))
			} else if v.operation == "history" {
				err = m.History()
			} else if v.operation == "wait" {
				time.Sleep(time.Duration(v.val.(int)) * time.Minute)
				return
			}
			assert.Nil(t, err)
			if v.operation == "history" {
				checkStringArrayHistory(v.output, ui.Lines, t)
			} else {
				checkStringArray(v.output, ui.Lines, t)
			}
		})
	}
}

func checkStringArray(expected []string, actual []string, t *testing.T) {
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], actual[i])
	}
}

func checkStringArrayHistory(expected []string, actual []string, t *testing.T) {
	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.True(t, strings.HasSuffix(actual[i], expected[i]))
	}
}
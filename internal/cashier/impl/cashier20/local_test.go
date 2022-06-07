package cashier20

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xhyyd/atm/internal/cashier"
	"testing"
)

// todo add more testcases
func TestImpl_Dispense(t *testing.T) {
	type Testcase struct {
		current float64
		val float64
		err error
		Expectation float64
	}
	type Testsuite []Testcase

	ts := Testsuite {
		{0,1.0,cashier.ErrNotEnoughCash,0},
		{2,3.0,cashier.ErrNotEnoughCash,0},
		{28,21.0,cashier.ErrInvalidRequest,0},
		{20,20,nil,0},
		{55,20,nil,35},
		{108.6,100.0,nil,8.6},
	}
	for i, v := range ts { //遍历
		l := NewCashierWithMoney(v.current)
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			err := l.Withdraw(v.val)
			if v.err != nil {
				assert.Equal(t, v.err, err)
			} else {
				assert.Equal(t, v.Expectation, l.Total())
			}
		})
	}
}

func TestImpl_WithdrawAdjust(t *testing.T) {
	type Testcase struct {
		current float64
		val float64
		err error
		Expectation float64
	}
	type Testsuite []Testcase

	ts := Testsuite {
		{0,1.0, errInvalidFormat,0},
		{2,3.0, errInvalidFormat,0},
		{28,21.0, errInvalidFormat,0},
		{20,20,nil,20},
		{108.6,100.0,nil,100},
		{85.6,100.0,nil,80},
	}
	for i, v := range ts { //遍历
		l := NewCashierWithMoney(v.current)
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			amount, err := l.WithdrawValidateAndAdjust(v.val)
			if v.err != nil {
				assert.Equal(t, v.err, err)
			} else {
				assert.Equal(t, v.Expectation, amount)
			}
		})
	}
}

func TestImpl_Deposit(t *testing.T) {
	type Testcase struct {
		current float64
		val float64
		Expectation float64
	}
	type Testsuite []Testcase

	ts := Testsuite {
		{0,1.0,1.0},
		{2,3.0,5},
		{28,21.0,49},
		{20,20,40},
		{108.6,100.0,208.6},
		{85.6,100.0,185.6},
	}
	for i, v := range ts { //遍历
		l := NewCashierWithMoney(v.current)
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			_ = l.Deposit(v.val)
			assert.Equal(t, v.Expectation, l.Total())
		})
	}
}

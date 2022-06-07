package local

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xhyyd/atm/internal/account"
	"testing"
)

func defaultManager() (account.Manager, error) {
	return NewManagerWithData([][]string{
		{"2859459814","7386","10.24"},
		{"1434597300","4557","90000.55"},
		{"7089382418","0075","0.00"},
		{"2001377812","5950","60.00"},
	})
}

func TestNewManagerWithData(t *testing.T) {
	type Testcase struct {
		data [][]string
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{ data: [][]string {
				{"2859459814","7386","10.24"},
				{"1434597300","4557","90000.55"},
				{"7089382418","0075","0.00"},
				{"2001377812","5950","60.00"},
			},
		}, {
			data: nil,
		},
	}
	for i, v := range ts { //遍历
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			_, err := NewManagerWithData(v.data)
			assert.Nil(t, err)
		})
	}
}

func TestManager_Authorize(t *testing.T) {
	type Testcase struct {
		id string
		pin string
		err error
		balance float64
	}
	type Testsuite []Testcase
	ts := Testsuite {
		{"2859459814","7386", nil, 10.24},
		{"1434597300","4557", nil, 90000.55},
		{"7089382418","7381", account.ErrInvalidPin, 0},
		{"7089382413","5950", account.ErrNotFound, 0},
	}

	for i, v := range ts { //遍历
		m, err := defaultManager()
		assert.Nil(t, err)
		t.Run(fmt.Sprintf("TestImpl_Dispense_%d" , i), func (t *testing.T){
			dl, err := m.Authorize(v.id, v.pin)
			if v.err != nil {
				assert.Equal(t, v.err, err)
			} else {
				assert.Nil(t, err)
				b, err := dl.Balance()
				assert.Nil(t, err)
				assert.Equal(t, v.balance, b)
			}
		})
	}
}

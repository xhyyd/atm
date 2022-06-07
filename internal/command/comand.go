package command

import (
	"github.com/pkg/errors"
	"github.com/xhyyd/atm/internal/atm"
	"strconv"
	"strings"
)

var EOF = errors.New("EOF")
var ErrInvalidCommand = errors.New("Invalid Command")

func Process(atm *atm.ATM, text string) error {
	tokens := strings.Split(text, " ")
	switch tokens[0] {
	case "":
		return nil
	case "authorize":
		if len(tokens) != 3 {
			return errors.New("Authorize: Invalidate input")
		}
		err := atm.Authorize(tokens[1], tokens[2])
		if err != nil {
			return errors.Wrap(err, "Authorize")
		}
	case "withdraw":
		if len(tokens) != 2 {
			return errors.New("Withdraw: Invalidate input")
		}
		val, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			return errors.Wrap(err, "Withdraw: Invalidate input")
		}
		err = atm.Withdraw(val)
		if err != nil {
			return errors.Wrap(err, "Withdraw")
		}
	case "deposit":
		if len(tokens) != 2 {
			return errors.New("Deposit: Invalidate input")
		}
		val, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			return errors.Wrap(err, "Deposit: Invalidate input")
		}
		err = atm.Deposit(val)
		if err != nil {
			return errors.Wrap(err, "Deposit")
		}
	case "balance":
		if len(tokens) != 1 {
			return errors.New("Balance: Invalidate input")
		}
		err := atm.Balance()
		if err != nil {
			return errors.Wrap(err, "Balance")
		}
	case "history":
		if len(tokens) != 1 {
			return errors.New("History: Invalidate input")
		}
		err := atm.History()
		if err != nil {
			return errors.Wrap(err, "History")
		}
	case "logout":
		if len(tokens) != 1 {
			return errors.New("Logout: Invalidate input")
		}
		atm.Logout()
	case "end":
		if len(tokens) != 1 {
			return errors.New("End: Invalidate input")
		}
		return EOF
	default:
		return ErrInvalidCommand
	}
	return nil
}


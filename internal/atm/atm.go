package atm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xhyyd/atm/internal/account"
	"github.com/xhyyd/atm/internal/cashier"
	"github.com/xhyyd/atm/internal/ui"
	"time"
)

const timeout = 2 * time.Minute

const AuthorizeRequired = "Authorization required."
const AuthorizeSuccess = "%s successfully authorized."
const AuthorizeFailed = "Authorization failed."
const Logout = "Account %s logged out"
const UnableToProcess = "Unable to process your withdrawal at this time."
const UnableToProcessFull = "Unable to dispense full amount requested at this time."
const AccountDispense = "Amount dispensed: $%.2f"
const CurrentBalance = "Current balance: %.2f"
const NoHistory = "No history found"
const NoAccountAuthorized = "No account is currently authorized"
const WithdrawWithFee = "You have been charged an overdraft fee of $%.2f. Current balance: %.2f"
const InvalidInput = "Invalid input, amount must be positive."
const Overdrawn = "Your account is overdrawn. Current balance: %.2f"

type ATM struct {
	accountManager account.Manager
	cashier        cashier.Cashier
	delegate       account.Delegate
	lastOperationTime time.Time
	ui                ui.UI
}

func NewATM(accountManager account.Manager, cashier cashier.Cashier, ui ui.UI) *ATM {
	return &ATM{
		accountManager: accountManager,
		cashier:        cashier,
		ui:             ui,
	}
}

func (atm *ATM) currentDelegate() account.Delegate {
	if time.Now().Sub(atm.lastOperationTime) > timeout {
		atm.delegate = nil
	}
	return atm.delegate
}

func (atm *ATM) updateOperationTime() {
	atm.lastOperationTime = time.Now()
}

func (atm *ATM) Authorize(id string, pin string) error {
	delegate, err := atm.accountManager.Authorize(id, pin)
	if err != nil {
		if err == account.ErrInvalidPin || err == account.ErrNotFound{
			atm.ui.Println(AuthorizeFailed)
			return nil
		}
		return err
	}
	atm.delegate = delegate
	atm.updateOperationTime()
	atm.ui.Println(fmt.Sprintf(AuthorizeSuccess, id))
	return nil
}

func (atm *ATM) Withdraw(value float64) error {
	// check account delegate is valid
	dl := atm.currentDelegate()
	if dl == nil {
		atm.ui.Println(AuthorizeRequired)
		return nil
	}

	// update operation time
	atm.updateOperationTime()

	// get real withdraw value from cashier, if not enough, adjust to the value in the machine
	realValue, err := atm.cashier.WithdrawValidateAndAdjust(value)
	if err != nil {
		atm.ui.Println(err.Error())
		return nil
	}
	if realValue == 0  {
		atm.ui.Println(UnableToProcess)
		return nil
	}
	if realValue != value {
		atm.ui.Println(UnableToProcessFull)
	}

	// withdraw from account delegate
	fee, balance, err := dl.Withdraw(realValue)
	if err != nil {
		if err == account.ErrOverdrawn {
			atm.ui.Println(fmt.Sprintf(Overdrawn, balance))
			return nil
		}
		// withdraw failed, this may happen in distribute environment
		return errors.Wrap(err, "account withdraw")
	}

	if fee == 0 {
		atm.ui.Println(fmt.Sprintf(AccountDispense, realValue))
		atm.ui.Println(fmt.Sprintf(CurrentBalance, balance))
	} else {
		atm.ui.Println(fmt.Sprintf(AccountDispense, realValue))
		atm.ui.Println(fmt.Sprintf(WithdrawWithFee, fee, balance))
	}

	// dispensing money from cashier
	err = atm.cashier.Withdraw(realValue)
	if err != nil {
		// dispense failed, this should not happen in local machine
		return errors.Wrap(err, "cashier dispense")
	}
	return nil
}

func (atm ATM) Deposit(value float64) error {
	// check account delegate is valid
	dl := atm.currentDelegate()
	if dl == nil {
		atm.ui.Println(AuthorizeRequired)
		return nil
	}

	// update operation time
	atm.updateOperationTime()

	// deposit with cashier
	err := atm.cashier.Deposit(value)
	if err != nil {
		if err == account.ErrInvalidAmount {
			atm.ui.Println(InvalidInput)
			return nil
		}
	}

	// deposit account delegate
	balance, err := dl.Deposit(value)
	if err != nil {
		return errors.Wrap(err, "account delegate deposit")
	}
	atm.ui.Println(fmt.Sprintf(CurrentBalance, balance))
	return nil
}

func (atm ATM) Balance() error {
	// check account delegate is valid
	dl := atm.currentDelegate()
	if dl == nil {
		atm.ui.Println(AuthorizeRequired)
		return nil
	}

	// update operation time
	atm.updateOperationTime()

	// get balance from account delegate
	balance, err := dl.Balance()
	if err != nil {
		return errors.Wrap(err, "account delegate get balance")
	}
	atm.ui.Println(fmt.Sprintf(CurrentBalance, balance))
	return nil
}

func (atm ATM) History() error {
	// check account delegate is valid
	dl := atm.currentDelegate()
	if dl == nil {
		atm.ui.Println(AuthorizeRequired)
		return nil
	}

	// update operation time
	atm.updateOperationTime()

	transactions, err := dl.History()
	if err != nil {
		return errors.Wrap(err, "account delegate get history")
	}

	if len(transactions) == 0 {
		atm.ui.Println(NoHistory)
		return nil
	}
	for _, trans := range transactions {
		atm.ui.Println(transToString(trans))
	}
	return nil
}

func transToString(trans account.Transaction) string {
	return fmt.Sprintf("%s %.2f %.2f", trans.Time.Format("2006-01-02 15:04:05"), trans.Value, trans.Balance)
}

func (atm ATM) Logout() {
	// check account delegate is valid
	dl := atm.currentDelegate()
	if dl == nil {
		atm.ui.Println(NoAccountAuthorized)
		return
	}

	id := dl.ID()
	atm.delegate = nil
	atm.ui.Println(fmt.Sprintf(Logout, id))
}
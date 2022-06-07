package account

import "github.com/pkg/errors"

var ErrNotFound = errors.New("Account not found")
var ErrInvalidPin = errors.New("Invalid pin")

/*
Manager is manager of the bank
Authorize do authentication and return delegate of the account
 */
type Manager interface {
	Authorize(id, pin string) (Delegate, error)
}
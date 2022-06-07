package account

import "time"

type Account struct {
	ID string	`json:"id"`
	Pin string `json:"pin"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	Time time.Time
	Value float64
	Balance float64
}
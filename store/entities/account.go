package entities

import (
	"time"
)

type Account struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewAccount(name string, balance float64) Account {
	now := time.Now()
	return Account{
		Name:      name,
		Balance:   balance,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

package entities

import (
	"time"

	decimal "github.com/shopspring/decimal"
)

type Account struct {
	Id        int             `db:"id"`
	Name      string          `db:"name"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func NewAccount(name string, balance decimal.Decimal) Account {
	now := time.Now()
	return Account{
		Name:      name,
		Balance:   balance,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

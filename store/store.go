package store

import (
	"context"
	"tiny-bank-api/pkg/database"
	"tiny-bank-api/store/entities"
)

type Store struct {
	db database.SQLDB
}

func NewStore(db database.SQLDB) Store {
	return Store{
		db: db,
	}
}

func (s Store) CreateAccount(ctx context.Context, name string, balance float64) error {
	account := entities.NewAccount(name, balance)
	q := `
		INSERT INTO accounts (name, balance, created_at, updated_at)
		VALUES (:name, :balance, :created_at, :updated_at);
	`
	_, err := s.db.NamedExecContext(ctx, q, account)
	return err
}

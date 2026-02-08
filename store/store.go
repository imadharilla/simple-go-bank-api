package store

import (
	"context"
	"log/slog"
	"tiny-bank-api/pkg/database"
	"tiny-bank-api/store/entities"

	"github.com/jmoiron/sqlx"
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

func (s Store) GetAccountById(ctx context.Context, accountId int64) (entities.Account, error) {
	var account entities.Account
	q := `SELECT id, name, balance, created_at, updated_at FROM accounts WHERE id = $1;`
	if err := s.db.QueryRowxContext(ctx, q, accountId).StructScan(&account); err != nil {
		return entities.Account{}, err
	}
	return account, nil
}

func (s Store) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	var accounts []entities.Account
	q := `SELECT id, name, balance, created_at, updated_at FROM accounts;`
	rows, err := s.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Warn("Failed to close rows", "error", err)
		}
	}()

	for rows.Next() {
		var account entities.Account
		if err := rows.StructScan(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s Store) AddBalance(ctx context.Context, accountId int64, amount float64) error {
	q := `
		UPDATE accounts 
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2;
	`
	_, err := s.db.ExecContext(ctx, q, amount, accountId)
	return err
}

func (s Store) AddBalanceWithTx(ctx context.Context, tx database.Querier, accountId int64, amount float64) error {
	q := `
		UPDATE accounts 
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2;
	`
	_, err := tx.ExecContext(ctx, q, amount, accountId)
	return err
}

func (s Store) SubtractBalanceWithTx(ctx context.Context, tx database.Querier, accountId int64, amount float64) error {
	q := `
		UPDATE accounts 
		SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2;
	`
	_, err := tx.ExecContext(ctx, q, amount, accountId)
	return err
}

func (s Store) GetAccountByIdWithTx(ctx context.Context, tx database.Querier, accountId int64) (entities.Account, error) {
	var account entities.Account
	q := `SELECT id, name, balance, created_at, updated_at FROM accounts WHERE id = $1;`
	if err := tx.QueryRowxContext(ctx, q, accountId).StructScan(&account); err != nil {
		return entities.Account{}, err
	}
	return account, nil
}

func (s Store) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, nil)
}

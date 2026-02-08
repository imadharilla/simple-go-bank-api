package database

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// Querier An interface to use for both sqlx.DB and sqlx.Tx (to use a transaction or not)
// We intentionally don't copy all methods so we can restrict our use cases to the methods
// that are best practice (i.e. always using context-based methods). We also intentionally
// exclude the `Must____` methods because we never need to panic
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type LoggingDB struct {
	SQLDB
	Logger *slog.Logger
}

var _ SQLDB = LoggingDB{}

func (q LoggingDB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	q.Logger.Debug(query)
	return q.SQLDB.QueryxContext(ctx, query, args...)
}

func (q LoggingDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	q.Logger.Debug(query)
	return q.SQLDB.QueryRowxContext(ctx, query, args...)
}

func (q LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	q.Logger.Debug(query)
	return q.SQLDB.ExecContext(ctx, query, args...)
}

func (q LoggingDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	q.Logger.Debug(query)
	return q.SQLDB.NamedExecContext(ctx, query, arg)
}

func (q LoggingDB) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	q.Logger.Debug(query)
	return q.SQLDB.PrepareNamedContext(ctx, query)
}

func (q LoggingDB) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	q.Logger.Debug(query)
	return q.SQLDB.PreparexContext(ctx, query)
}

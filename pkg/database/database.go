package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

type SQLDB interface {
	Querier
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
	Close() error
}

// NewConnection constructor for database connection pool to postgres
func NewConnection(ctx context.Context, connString string) (*sqlx.DB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	stdDB, err := otelsql.Open("pgx", config.ConnString(), otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
	if err != nil {
		return nil, fmt.Errorf("error connecting to pg: %w", err)
	}

	db := sqlx.NewDb(stdDB, "pgx")

	return db, nil
}

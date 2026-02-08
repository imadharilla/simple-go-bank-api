package store

import (
	"tiny-bank-api/pkg/database"
)

type Store struct {
	db database.SQLDB
}

func NewStore(db database.SQLDB) Store {
	return Store{
		db: db,
	}
}

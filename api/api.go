package api

import (
	"log/slog"
	"tiny-bank-api/store"
)

type API struct {
	logger *slog.Logger
	store  store.Store
}

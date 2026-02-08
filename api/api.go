package api

import (
	"context"
	"fmt"
	"log/slog"
	"tiny-bank-api/store"
)

type API struct {
	logger *slog.Logger
	store  store.Store
}

func NewAPI(logger *slog.Logger, store store.Store) *API {
	return &API{
		logger: logger,
		store:  store,
	}
}

func (s API) GetAccounts(ctx context.Context, request GetAccountsRequestObject) (GetAccountsResponseObject, error) {
	//TODO implement me
	return nil, fmt.Errorf("not implemented yet")
}

func (s API) CreateAccount(ctx context.Context, request CreateAccountRequestObject) (CreateAccountResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

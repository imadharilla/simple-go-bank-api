package api

import (
	"context"
	"database/sql"
	"errors"
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
	accounts, err := s.store.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	response := make(GetAccounts200JSONResponse, 0, len(accounts))
	for _, acc := range accounts {
		response = append(response, Account{
			Id:        int64(acc.Id),
			Name:      acc.Name,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		})
	}

	return response, nil
}

func (s API) CreateAccount(ctx context.Context, request CreateAccountRequestObject) (CreateAccountResponseObject, error) {
	err := s.store.CreateAccount(ctx, request.Body.Name, 0)
	if err != nil {
		return nil, err
	}
	return CreateAccount201Response{}, nil
}

func (s API) AddBalanceToAccount(ctx context.Context, request AddBalanceToAccountRequestObject) (AddBalanceToAccountResponseObject, error) {
	if request.Body.Amount <= 0 {
		return AddBalanceToAccount400Response{}, nil
	}

	// check if the account exists
	_, err := s.store.GetAccountById(ctx, request.AccountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AddBalanceToAccount404Response{}, err
		}
		return nil, err
	}

	err = s.store.AddBalance(ctx, request.AccountId, request.Body.Amount)
	if err != nil {
		return nil, err
	}

	return AddBalanceToAccount200Response{}, nil
}

func (s API) TransferMoney(ctx context.Context, request TransferMoneyRequestObject) (TransferMoneyResponseObject, error) {
	if request.Body.Amount <= 0 {
		return TransferMoney400JSONResponse{Message: "amount must be greater than 0"}, nil
	}
	if request.AccountId == request.Body.TargetAccountId {
		return TransferMoney400JSONResponse{Message: "cannot transfer to the same account"}, nil
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			slog.Warn("Failed to rollback the transaction", "error", err)
		}
	}()

	// check target account exists
	_, err = s.store.GetAccountByIdWithTx(ctx, tx, request.Body.TargetAccountId)
	if err != nil {
		return TransferMoney400JSONResponse{Message: "target account not found"}, nil
	}

	// Check source account exists and has sufficient balance
	sourceAccount, err := s.store.GetAccountByIdWithTx(ctx, tx, request.AccountId)
	if err != nil {
		return TransferMoney400JSONResponse{Message: "source account not found"}, nil
	}
	if sourceAccount.Balance < request.Body.Amount {
		return TransferMoney400JSONResponse{Message: "insufficient balance"}, nil
	}

	err = s.store.SubtractBalanceWithTx(ctx, tx, request.AccountId, request.Body.Amount)
	if err != nil {
		return nil, err
	}
	err = s.store.AddBalanceWithTx(ctx, tx, request.Body.TargetAccountId, request.Body.Amount)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return TransferMoney200Response{}, nil
}

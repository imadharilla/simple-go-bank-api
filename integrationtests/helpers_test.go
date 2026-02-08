package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiny-bank-api/api"
)

func requireStatus(t *testing.T, expected int, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != expected {
		t.Fatalf("expected status %d, got %d. Body: %s", expected, rec.Code, rec.Body.String())
	}
}

func reqGETAccounts(t *testing.T, handler http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func reqPOSTAccount(t *testing.T, handler http.Handler, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func mustPOSTAccount(t *testing.T, handler http.Handler, name string) {
	t.Helper()
	rec := reqPOSTAccount(t, handler, map[string]any{
		"name": name,
	})
	requireStatus(t, http.StatusCreated, rec)
}

func mustGETAccounts(t *testing.T, handler http.Handler) []api.Account {
	t.Helper()
	rec := reqGETAccounts(t, handler)
	requireStatus(t, http.StatusOK, rec)

	var accounts []api.Account
	if err := json.NewDecoder(rec.Body).Decode(&accounts); err != nil {
		t.Fatalf("failed to decode accounts response: %v", err)
	}
	return accounts
}

func requireAccountExists(t *testing.T, handler http.Handler, name string) api.Account {
	t.Helper()
	accounts := mustGETAccounts(t, handler)
	for _, acc := range accounts {
		if acc.Name == name {
			return acc
		}
	}
	t.Fatalf("account with name %q not found", name)
	return api.Account{} // unreachable but required for compilation
}

func reqPOSTAddBalance(t *testing.T, handler http.Handler, accountId int64, amount float64) *httptest.ResponseRecorder {
	t.Helper()
	jsonBody, err := json.Marshal(map[string]any{
		"amount": amount,
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/accounts/%d/add-balance", accountId), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func mustPOSTAddBalance(t *testing.T, handler http.Handler, accountId int64, amount float64) {
	t.Helper()
	rec := reqPOSTAddBalance(t, handler, accountId, amount)
	requireStatus(t, http.StatusOK, rec)
}

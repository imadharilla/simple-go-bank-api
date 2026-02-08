package integrationtests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiny-bank-api/api"
)

// reqGETAccounts performs a GET request to /api/accounts and returns the response recorder
func reqGETAccounts(t *testing.T, handler http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// reqPOSTAccount performs a POST request to /api/accounts with the given body and returns the response recorder
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

// mustPOSTAccount creates a new account and asserts the response status is 201 Created
func mustPOSTAccount(t *testing.T, handler http.Handler, name string) {
	t.Helper()
	rec := reqPOSTAccount(t, handler, map[string]any{
		"name": name,
	})
	requireStatus(t, http.StatusCreated, rec)
}

// mustGETAccounts retrieves all accounts and asserts the response status is 200 OK,
// returning the parsed list of accounts
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

// generateAccounts creates multiple accounts with the given names
func generateAccounts(t *testing.T, handler http.Handler, names []string) {
	t.Helper()
	for _, name := range names {
		mustPOSTAccount(t, handler, name)
	}
}

// requireStatus asserts that the response recorder has the expected status code
func requireStatus(t *testing.T, expected int, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != expected {
		t.Fatalf("expected status %d, got %d. Body: %s", expected, rec.Code, rec.Body.String())
	}
}

// requireAccountExists checks if an account with the given name exists
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

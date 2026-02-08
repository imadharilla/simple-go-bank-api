package integrationtests

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

var testHandler http.Handler

func TestAccountCreation(t *testing.T) {
	accountName := fmt.Sprintf("Aimad Creator - %d", time.Now().Unix())
	mustPOSTAccount(t, testHandler, accountName)
	requireAccountExists(t, testHandler, accountName)
}

func TestAddBalance(t *testing.T) {
	t.Run(`should fail if account doesn't exist`, func(t *testing.T) {
		rec := reqPOSTAddBalance(t, testHandler, 20321, 13.37)
		requireStatus(t, http.StatusInternalServerError, rec)
	})

	t.Run(`should fail if amount is zero or negative`, func(t *testing.T) {
		accountName := fmt.Sprintf("Aimad Negative Balance - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, accountName)
		account := requireAccountExists(t, testHandler, accountName)

		rec := reqPOSTAddBalance(t, testHandler, account.Id, -10)
		requireStatus(t, http.StatusBadRequest, rec)

		rec = reqPOSTAddBalance(t, testHandler, account.Id, 0)
		requireStatus(t, http.StatusBadRequest, rec)
	})

	t.Run(`should add balance successfully`, func(t *testing.T) {
		accountName := fmt.Sprintf("Aimad Add Balance - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, accountName)
		account := requireAccountExists(t, testHandler, accountName)

		// add 50
		mustPOSTAddBalance(t, testHandler, account.Id, 50)
		// add another 3.34
		mustPOSTAddBalance(t, testHandler, account.Id, 3.34)

		updatedAccount := requireAccountExists(t, testHandler, accountName)
		if updatedAccount.Balance != account.Balance+50+3.34 {
			t.Fatalf("expected balance to be %.2f but got %.2f", account.Balance+50, updatedAccount.Balance)
		}
	})

}

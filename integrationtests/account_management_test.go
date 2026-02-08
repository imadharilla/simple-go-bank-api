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
		requireStatus(t, http.StatusNotFound, rec)
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

func TestTransferMoney(t *testing.T) {
	t.Run(`should fail if source account doesn't exist`, func(t *testing.T) {
		targetName := fmt.Sprintf("Transfer Target 1 - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, targetName)
		targetAccount := requireAccountExists(t, testHandler, targetName)

		rec := reqPOSTTransfer(t, testHandler, 99999, targetAccount.Id, 10)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "source account not found", rec)
	})

	t.Run(`should fail if target account doesn't exist`, func(t *testing.T) {
		sourceName := fmt.Sprintf("Transfer Source 1 - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, sourceName)
		sourceAccount := requireAccountExists(t, testHandler, sourceName)
		mustPOSTAddBalance(t, testHandler, sourceAccount.Id, 100)

		rec := reqPOSTTransfer(t, testHandler, sourceAccount.Id, 99999, 10)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "target account not found", rec)
	})

	t.Run(`should fail if amount is zero or negative`, func(t *testing.T) {
		sourceName := fmt.Sprintf("Transfer Source 2 - %d", time.Now().Unix())
		targetName := fmt.Sprintf("Transfer Target 2 - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, sourceName)
		mustPOSTAccount(t, testHandler, targetName)
		sourceAccount := requireAccountExists(t, testHandler, sourceName)
		targetAccount := requireAccountExists(t, testHandler, targetName)

		rec := reqPOSTTransfer(t, testHandler, sourceAccount.Id, targetAccount.Id, 0)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "amount must be greater than 0", rec)

		rec = reqPOSTTransfer(t, testHandler, sourceAccount.Id, targetAccount.Id, -10)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "amount must be greater than 0", rec)
	})

	t.Run(`should fail if transferring to the same account`, func(t *testing.T) {
		sourceName := fmt.Sprintf("Transfer Same Account - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, sourceName)
		sourceAccount := requireAccountExists(t, testHandler, sourceName)

		rec := reqPOSTTransfer(t, testHandler, sourceAccount.Id, sourceAccount.Id, 10)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "cannot transfer to the same account", rec)
	})

	t.Run(`should fail if insufficient balance`, func(t *testing.T) {
		sourceName := fmt.Sprintf("Transfer Source 3 - %d", time.Now().Unix())
		targetName := fmt.Sprintf("Transfer Target 3 - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, sourceName)
		mustPOSTAccount(t, testHandler, targetName)
		sourceAccount := requireAccountExists(t, testHandler, sourceName)
		targetAccount := requireAccountExists(t, testHandler, targetName)
		mustPOSTAddBalance(t, testHandler, sourceAccount.Id, 50)

		rec := reqPOSTTransfer(t, testHandler, sourceAccount.Id, targetAccount.Id, 100)
		requireStatus(t, http.StatusBadRequest, rec)
		requireErrorMessage(t, "insufficient balance", rec)
	})

	t.Run(`should transfer money successfully`, func(t *testing.T) {
		sourceName := fmt.Sprintf("Transfer Source 4 - %d", time.Now().Unix())
		targetName := fmt.Sprintf("Transfer Target 4 - %d", time.Now().Unix())
		mustPOSTAccount(t, testHandler, sourceName)
		mustPOSTAccount(t, testHandler, targetName)
		sourceAccount := requireAccountExists(t, testHandler, sourceName)
		targetAccount := requireAccountExists(t, testHandler, targetName)

		mustPOSTAddBalance(t, testHandler, sourceAccount.Id, 100)
		mustPOSTTransfer(t, testHandler, sourceAccount.Id, targetAccount.Id, 30)

		updatedSource := requireAccountExists(t, testHandler, sourceName)
		updatedTarget := requireAccountExists(t, testHandler, targetName)

		if updatedSource.Balance != 70 {
			t.Fatalf("expected source balance to be 70 but got %.2f", updatedSource.Balance)
		}
		if updatedTarget.Balance != 30 {
			t.Fatalf("expected target balance to be 30 but got %.2f", updatedTarget.Balance)
		}
	})
}

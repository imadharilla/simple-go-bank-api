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

package sdk

import (
	"testing"
)

func TestVaultLogin(t *testing.T) {
	vault := NewVault()
	err := vault.Login("email", "password")
	checkTestError(t, err, "Could not log into vault")
}

func TestVaultStore(t *testing.T) {
	vault := NewVault()
	err := vault.Login("email", "password")
	checkTestError(t, err, "Could not log into vault")
	inputValue := "test value"
	err = vault.Put("/test", inputValue)
	checkTestError(t, err, "Could not store data")
	returnedValue, err := vault.Get("/test")
	checkTestError(t, err, "Could not retrieve data")
	if returnedValue != returnedValue {
		t.Fatalf("Returned value did not match: (%s) -> (%s)", inputValue, returnedValue)
	}
}

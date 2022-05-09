package main

import (
	"testing"
)

func TestVaultLogin(t *testing.T) {
	vault := newVault()
	err := vault.login("email", "password")
	if err != nil {
		t.Fatalf("Could not log into vault: %s", err)
	}
}

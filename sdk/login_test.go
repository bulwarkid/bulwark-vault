package main

import (
	"testing"
)

func TestLoginSecret(t *testing.T) {
	loginSecret, err := deriveLoginSecret("email", "password")
	checkTestError(t, err, "Failed to derive login secret")
	if len(loginSecret) != 32 {
		t.Fatal("Login secret wrong size")
	}
}

func TestGetMasterSecret(t *testing.T) {
	masterSecret := testMasterSecret(t)
	if len(masterSecret) != 32 {
		t.Fatal("Master secret invalid:", masterSecret)
	}
}

package main

import (
	"testing"
)

func TestLoginSecret(t *testing.T) {
	_, err := DeriveLoginSecret("email", "password")
	if err != nil {
		t.Log("Login secret error:", err)
		t.Fail()
	}
}

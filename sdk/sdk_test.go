package main

import (
	"reflect"
	"testing"
)

func TestLoginSecret(t *testing.T) {
	loginSecret, err := DeriveLoginSecret("email", "password")
	if err != nil {
		t.Log("Login secret error:", err)
		t.FailNow()
	}
	if len(loginSecret) != 32 {
		t.Log("Login secret wrong size")
		t.FailNow()
	}
}

func TestGetSalt(t *testing.T) {
	saltId := SaltIdForEmail("email")
	salt, err := GetSalt(saltId)
	if err != nil {
		t.Log("Failed to get salt:", err)
		t.FailNow()
	}
	if len(salt) != 32 {
		t.Log("Invalid salt:", salt, err)
		t.FailNow()
	}
	salt2, err := GetSalt(saltId)
	if err != nil {
		t.Log("Failed to get second salt:", err)
		t.FailNow()
	}
	if len(salt2) != 32 {
		t.Log("Invalid salt:", salt, err)
		t.FailNow()
	}
	if !reflect.DeepEqual(salt, salt2) {
		t.Log("Salt is not saved between fetches:", salt, salt2)
		t.FailNow()
	}
}

func TestGetMasterSecret(t *testing.T) {
	loginSecret, err := DeriveLoginSecret("email", "password")
	if err != nil {
		t.Log("Could not get login secret:", err)
		t.FailNow()
	}
	masterSecret, err := GetMasterSecret(loginSecret)
	if err != nil {
		t.Log("Could not get master secret:", err)
		t.FailNow()
	}
	if len(masterSecret) != 32 {
		t.Log("Master secret invalid:", masterSecret)
		t.FailNow()
	}
}

package main

import (
	"reflect"
	"testing"
)

func TestLoginSecret(t *testing.T) {
	loginSecret, err := deriveLoginSecret("email", "password")
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
	salt, err := getSalt("email")
	if err != nil {
		t.Log("Failed to get salt:", err)
		t.FailNow()
	}
	if len(salt) != 32 {
		t.Log("Invalid salt:", salt, err)
		t.FailNow()
	}
	salt2, err := getSalt("email")
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
	loginSecret, err := deriveLoginSecret("email", "password")
	if err != nil {
		t.Log("Could not get login secret in test:", err)
		t.FailNow()
	}
	masterSecret, err := getMasterSecret(loginSecret)
	if err != nil {
		t.Log("Could not get master secret in test:", err)
		t.FailNow()
	}
	if len(masterSecret) != 32 {
		t.Log("Master secret invalid:", masterSecret)
		t.FailNow()
	}
}

func TestGetKeyDirectory(t *testing.T) {
	masterSecret, err := newMasterSecret()
	if err != nil {
		t.Logf("Error generating master secret: %s", err)
		t.FailNow()
	}
	directory, err := getKeyDirectory(masterSecret)
	if err != nil {
		t.Logf("Error getting key directory: %s", err)
		t.FailNow()
	}
	var data AccessData
	data.accessKey = "access"
	data.encryptionKey = "encryption"
	(*directory)["/test"] = &data
	err = directory.Store(masterSecret)
	if err != nil {
		t.Logf("Error storing directory: %s", err)
		t.FailNow()
	}
}

func TestVaultLogin(t *testing.T) {
	vault := newVault()
	err := vault.login("email", "password")
	if err != nil {
		t.Logf("Could not log into vault: %s", err)
		t.FailNow()
	}
}

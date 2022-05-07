package main

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"testing"
)

func TestLoginSecret(t *testing.T) {
	loginSecret, err := DeriveLoginSecret("email", "password")
	if err != nil {
		t.Log("Login secret error:", err)
		t.FailNow()
	}
	if len(loginSecret) != 44 {
		t.Log("Login secret wrong size")
		t.FailNow()
	}
	_, err = base64.URLEncoding.DecodeString(loginSecret)
	if err != nil {
		t.Log("Login secret not base64")
		t.FailNow()
	}
}

func TestGetSalt(t *testing.T) {
	data := make([]byte, 32)
	rand.Read(data)
	saltIdBytes := hashSha256([]byte(data))
	saltId := base64.URLEncoding.EncodeToString(saltIdBytes)
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

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func deriveLoginSecret(email string, password string) ([]byte, error) {
	salt, err := getSalt(email)
	if err != nil {
		return nil, err
	}
	return bytesFromLowEntropy("login_secret:"+email+":"+password, salt), nil
}

func GetMasterSecret(loginSecret []byte) ([]byte, error) {
	secretBase64, err := getObjectByPath(loginSecret, "/master-secret")
	if err != nil {
		if !isReturnCode(err, 404) {
			return nil, fmt.Errorf("Could not get master secret: %w", err)
		}
		secretBytes, err := NewMasterSecret()
		if err != nil {
			return nil, fmt.Errorf("Could not generate random bytes: %w", err)
		}
		secretBase64 = base64.URLEncoding.EncodeToString(secretBytes)
		if err = writeObjectByPath(loginSecret, "/master-secret", secretBase64); err != nil {
			return nil, fmt.Errorf("Could not write master secret: %w", err)
		}
		return secretBytes, nil
	}
	masterSecret, err := base64.URLEncoding.DecodeString(secretBase64)
	if err != nil {
		return nil, fmt.Errorf("Invalid master secret returned from vault: %w", err)
	}
	return masterSecret, nil
}

func NewMasterSecret() ([]byte, error) {
	secretBytes := make([]byte, 32)
	n, err := rand.Read(secretBytes)
	if err != nil || n != 32 {
		return nil, fmt.Errorf("Could not generate random bytes: %w", err)
	}
	return secretBytes, nil
}

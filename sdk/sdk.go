package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

func get(url string) (string, error) {
	response, err := http.Get("http://localhost:8080" + url)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func hashSha256(data []byte) []byte {
	bytes := sha256.Sum256(data)
	return bytes[:]
}

func decryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptedData := make([]byte, len(data))
	block.Encrypt(encryptedData, data)
	return encryptedData, nil
}

func deterministicObjectAccessKey(secret string, path string) (string, error) {
	inputData := "object_access_key:" + secret + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func deterministicObjectEncryptionKey(secret string, path string) ([]byte, error) {
	inputData := "object_encryption_key:" + secret + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GetSalt(saltId string) ([]byte, error) {
	saltBase64, err := get("/vault/salt/" + saltId)
	if err != nil {
		return nil, err
	}
	salt, err := base64.URLEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GetObject(secret string, path string) (string, error) {
	accessKey, err := deterministicObjectAccessKey(secret, path)
	if err != nil {
		return "", err
	}
	objectBase64, err := get("/vault/object/" + accessKey)
	if err != nil {
		return "", err
	}
	encryptedBytes, err := base64.URLEncoding.DecodeString(objectBase64)
	if err != nil {
		return "", err
	}
	encryptionKey, err := deterministicObjectEncryptionKey(secret, path)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := decryptBytes(encryptedBytes, encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

func DeriveLoginSecret(email string, password string) (string, error) {
	saltId := base64.URLEncoding.EncodeToString(hashSha256([]byte(email)))
	salt, err := GetSalt(saltId)
	if err != nil {
		return "", err
	}
	inputData := "login_secret:" + email + ":" + password
	keyBytes := pbkdf2.Key([]byte(inputData), salt, 10000, 32, sha256.New)
	key := base64.URLEncoding.EncodeToString(keyBytes)
	return key, nil
}

func GetMasterSecret(loginSecret string) (string, error) {
	return GetObject(loginSecret, "/master-secret")
}

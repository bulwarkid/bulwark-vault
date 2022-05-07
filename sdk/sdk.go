package main

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

const BASE_URL = "http://localhost:8080"

type UnsuccessfulRequestError struct {
	code     int
	response string
}

func (e *UnsuccessfulRequestError) Error() string {
	return fmt.Sprintf("Unsuccessful request (%d): %s", e.code, e.response)
}

func get(path string) (string, error) {
	response, err := http.Get(BASE_URL + path)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("%w", &UnsuccessfulRequestError{code: response.StatusCode, response: string(bytes)})
	}
	return string(bytes), nil
}

func post(path string, dataType string, data string) (string, error) {
	response, err := http.Post(BASE_URL+path, dataType, strings.NewReader(data))
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
	decryptedData := make([]byte, len(data))
	block.Decrypt(decryptedData, data)
	return decryptedData, nil
}

func encryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptedData := make([]byte, len(data))
	block.Encrypt(encryptedData, data)
	return encryptedData, nil
}

func deterministicObjectAccessKey(secret []byte, path string) (string, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	inputData := "object_access_key:" + secretString + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func deterministicObjectEncryptionKey(secret []byte, path string) ([]byte, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	inputData := "object_encryption_key:" + secretString + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func SaltIdForEmail(email string) string {
	saltIdBytes := hashSha256([]byte(email))
	return base64.URLEncoding.EncodeToString(saltIdBytes)
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

func GetObjectByPath(secret []byte, path string) (string, error) {
	accessKey, err := deterministicObjectAccessKey(secret, path)
	if err != nil {
		return "", fmt.Errorf("Could not get access key: %w", err)
	}
	encryptionKey, err := deterministicObjectEncryptionKey(secret, path)
	if err != nil {
		return "", fmt.Errorf("Could not get encryption key: %w", err)
	}
	return GetObject(accessKey, encryptionKey)
}

func GetObject(accessKey string, encryptionKey []byte) (string, error) {
	objectBase64, err := get("/vault/object/" + accessKey)
	if err != nil {
		return "", fmt.Errorf("Could not retrieve object data: %w", err)
	}
	encryptedBytes, err := base64.URLEncoding.DecodeString(objectBase64)
	if err != nil {
		return "", fmt.Errorf("Could not decode object bytes: %w (%s)", err, objectBase64)
	}
	decryptedBytes, err := decryptBytes(encryptedBytes, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("Could not decrypt bytes: %w", err)
	}
	return string(decryptedBytes), nil
}

func WriteObjectByPath(secret []byte, path string, data string) error {
	accessKey, err := deterministicObjectAccessKey(secret, path)
	if err != nil {
		return err
	}
	encryptionKey, err := deterministicObjectEncryptionKey(secret, path)
	if err != nil {
		return err
	}
	return WriteObject(accessKey, encryptionKey, data)
}

func WriteObject(accessKey string, encryptionKey []byte, data string) error {
	encryptedBytes, err := encryptBytes([]byte(data), encryptionKey)
	if err != nil {
		return err
	}
	objectBase64 := base64.URLEncoding.EncodeToString(encryptedBytes)
	_, err = post("/vault/object/"+accessKey, "text/plain", objectBase64)
	if err != nil {
		return err
	}
	return nil
}

func DeriveLoginSecret(email string, password string) ([]byte, error) {
	saltId := base64.URLEncoding.EncodeToString(hashSha256([]byte(email)))
	salt, err := GetSalt(saltId)
	if err != nil {
		return nil, err
	}
	inputData := "login_secret:" + email + ":" + password
	return pbkdf2.Key([]byte(inputData), salt, 10000, 32, sha256.New), nil
}

func GetMasterSecret(loginSecret []byte) ([]byte, error) {
	secretBase64, err := GetObjectByPath(loginSecret, "/master-secret")
	if err != nil {
		var requestError *UnsuccessfulRequestError
		if !errors.As(err, &requestError) {
			return nil, fmt.Errorf("Could not get master secret: %w", err)
		}
		if requestError.code != 404 {
			return nil, fmt.Errorf("Could not get master secret: %w", err)
		}
		secretBytes := make([]byte, 32)
		_, err = rand.Read(secretBytes)
		if err != nil {
			return nil, fmt.Errorf("Could not generate random bytes: %w", err)
		}
		secretBase64 = base64.URLEncoding.EncodeToString(secretBytes)
		if err = WriteObjectByPath(loginSecret, "/master-secret", secretBase64); err != nil {
			return nil, fmt.Errorf("Could not write master secret: %w", err)
		}
		return secretBytes, nil
	}
	return base64.URLEncoding.DecodeString(secretBase64)
}

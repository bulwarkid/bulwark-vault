package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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

func isReturnCode(err error, code int) bool {
	var requestError *UnsuccessfulRequestError
	return errors.As(err, &requestError) && requestError.code == code
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

func randomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	n, err := rand.Read(bytes)
	if err != nil || n != length {
		return nil, err
	}
	return bytes, nil
}

func decryptBytes(data []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	decryptedData, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func encryptBytes(data []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	nonce, err := randomBytes(12)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	encryptedData := gcm.Seal(nil, nonce, data, nil)
	return encryptedData, nonce, nil
}

func deterministicObjectAccessKey(secret []byte, path string) (string, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	inputData := "object_access_key:" + secretString + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return "", fmt.Errorf("Could not generate access key: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func deterministicObjectEncryptionKey(secret []byte, path string) ([]byte, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	inputData := "object_encryption_key:" + secretString + ":" + path
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, 32))
	if err != nil {
		return nil, fmt.Errorf("Could not generate encryption key: %w", err)
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

type EncryptedBlob struct {
	Blob string `json:"blob"`
	Iv   string `json:"iv"`
}

func encodeBlob(blob []byte, iv []byte) (string, error) {
	blobBase64 := base64.URLEncoding.EncodeToString(blob)
	ivBase64 := base64.URLEncoding.EncodeToString(iv)
	blobJson, err := json.Marshal(&EncryptedBlob{Blob: blobBase64, Iv: ivBase64})
	if err != nil {
		return "", err
	}
	return string(blobJson), nil
}

func decodeBlob(blobJson string) ([]byte, []byte, error) {
	var blob EncryptedBlob
	if err := json.Unmarshal([]byte(blobJson), &blob); err != nil {
		return nil, nil, err
	}
	blobBytes, err := base64.URLEncoding.DecodeString(blob.Blob)
	if err != nil {
		return nil, nil, err
	}
	ivBytes, err := base64.URLEncoding.DecodeString(blob.Iv)
	if err != nil {
		return nil, nil, err
	}
	return blobBytes, ivBytes, nil
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
	blobJson, err := base64.URLEncoding.DecodeString(objectBase64)
	if err != nil {
		return "", fmt.Errorf("Could not decode object bytes: %w (%s)", err, objectBase64)
	}
	blob, iv, err := decodeBlob(string(blobJson))
	if err != nil {
		return "", fmt.Errorf("Could not decode blob json: %w (%s)", err, blobJson)
	}
	decryptedBytes, err := decryptBytes(blob, encryptionKey, iv)
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
	encryptedBytes, iv, err := encryptBytes([]byte(data), encryptionKey)
	if err != nil {
		return fmt.Errorf("Could not encrypt data: %w", err)
	}
	blobJson, err := encodeBlob(encryptedBytes, iv)
	if err != nil {
		return fmt.Errorf("Could not create json: %w", err)
	}
	objectBase64 := base64.URLEncoding.EncodeToString([]byte(blobJson))
	_, err = post("/vault/object/"+accessKey, "text/plain", objectBase64)
	if err != nil {
		return fmt.Errorf("Could not store encrypted data: %w", err)
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
		if !isReturnCode(err, 404) {
			return nil, fmt.Errorf("Could not get master secret: %w", err)
		}
		secretBytes, err := NewMasterSecret()
		if err != nil {
			return nil, fmt.Errorf("Could not generate random bytes: %w", err)
		}
		secretBase64 = base64.URLEncoding.EncodeToString(secretBytes)
		if err = WriteObjectByPath(loginSecret, "/master-secret", secretBase64); err != nil {
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

type AccessData struct {
	accessKey     string
	encryptionKey string
}

type KeyDirectory map[string]*AccessData

func NewKeyDirectory() *KeyDirectory {
	directory := make(KeyDirectory)
	return &directory
}

func (directory *KeyDirectory) Load(masterSecret []byte) error {
	keyDirectoryObject, err := GetObjectByPath(masterSecret, "/directory")
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(keyDirectoryObject), directory); err != nil {
		return err
	}
	return nil
}

func (directory *KeyDirectory) Store(masterSecret []byte) error {
	jsonData, err := json.Marshal(directory)
	if err != nil {
		return err
	}
	return WriteObjectByPath([]byte(jsonData), "/master-secret", string(jsonData))
}

func getKeyDirectory(masterSecret []byte) (*KeyDirectory, error) {
	directory := NewKeyDirectory()
	err := directory.Load(masterSecret)
	if err != nil {
		if !isReturnCode(err, 404) {
			return nil, err
		}
		if err = directory.Store(masterSecret); err != nil {
			return nil, err
		}
	}
	return directory, nil
}

type VaultAccess struct {
	masterSecret []byte
	keyDirectory *KeyDirectory
}

func NewVault() *VaultAccess {
	return &VaultAccess{masterSecret: nil, keyDirectory: nil}
}

func (access *VaultAccess) Login(email, password string) error {
	loginSecret, err := DeriveLoginSecret(email, password)
	if err != nil {
		return err
	}
	masterSecret, err := GetMasterSecret(loginSecret)
	if err != nil {
		return err
	}
	access.masterSecret = masterSecret
	keyDirectory, err := getKeyDirectory(access.masterSecret)
	if err != nil {
		return err
	}
	access.keyDirectory = keyDirectory
	return nil
}

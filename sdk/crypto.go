package sdk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

const ENCRYPTION_KEY_LENGTH = 32

type AESEncryptionKey []byte

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

func randomEncryptionKey() (AESEncryptionKey, error) {
	bytes, err := randomBytes(ENCRYPTION_KEY_LENGTH)
	if err != nil {
		return nil, fmt.Errorf("Could not generate encryption key: %w", err)
	}
	return AESEncryptionKey(bytes), nil
}

func deterministicObjectEncryptionKey(secret []byte, path string) (AESEncryptionKey, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	bytes, err := bytesFromHighEntropy("object_encryption_key:"+secretString+":"+path, ENCRYPTION_KEY_LENGTH)
	if err != nil {
		return nil, fmt.Errorf("Could not generate encryption key: %w", err)
	}
	return AESEncryptionKey(bytes), nil
}

func decryptBytes(data []byte, key AESEncryptionKey, nonce []byte) ([]byte, error) {
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

func encryptBytes(data []byte, key AESEncryptionKey) ([]byte, []byte, error) {
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

func bytesFromLowEntropy(inputData string, salt []byte, length int) []byte {
	return pbkdf2.Key([]byte(inputData), salt, 10000, length, sha256.New)
}

func bytesFromHighEntropy(inputData string, length int64) ([]byte, error) {
	r := hkdf.New(sha256.New, []byte(inputData), nil, nil)
	bytes, err := ioutil.ReadAll(io.LimitReader(r, length))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type PublicKeyPair struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func randomPublicKeyPair() (*PublicKeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("Could not generate ed25519 key pair: %w", err)
	}
	return &PublicKeyPair{publicKey: publicKey, privateKey: privateKey}, nil
}

func signData(privateKey ed25519.PrivateKey, data []byte) []byte {
	return ed25519.Sign(privateKey, data)
}

func verifySignature(publicKey ed25519.PublicKey, data []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}

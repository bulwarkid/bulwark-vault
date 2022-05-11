package sdk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

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

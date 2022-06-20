package sdk

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
)

type AuthenticatedDataJson struct {
	Data      string `json:"data"`
	Iv        string `json:"iv"`
	Signature string `json:"signature"`
}

func getAuthData(publicKey ed25519.PublicKey, encryptionKey AESEncryptionKey) ([]byte, error) {
	jsonData, err := get("/vault/authenticated_object/" + b64encode(publicKey))
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve data: %w", err)
	}
	dataBlob := AuthenticatedDataJson{}
	err = json.Unmarshal([]byte(jsonData), &dataBlob)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	data := b64decode(dataBlob.Data)
	iv := b64decode(dataBlob.Iv)
	signature := b64decode(dataBlob.Signature)
	if !verifySignature(publicKey, data, signature) {
		return nil, fmt.Errorf("Authentication failed for data for public key %s", b64encode(publicKey))
	}
	decryptedData, err := decryptBytes(data, encryptionKey, iv)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting bytes: %w", err)
	}
	return decryptedData, nil
}

func writeAuthData(data []byte, keyPair *PublicKeyPair, encryptionKey AESEncryptionKey) error {
	encryptedData, nonce, err := encryptBytes(data, encryptionKey)
	if err != nil {
		return fmt.Errorf("Could not encrypt bytes: %w", err)
	}
	signature := signData(keyPair.privateKey, encryptedData)
	jsonData, err := json.MarshalIndent(
		AuthenticatedDataJson{
			Data:      b64encode(encryptedData),
			Iv:        b64encode(nonce),
			Signature: b64encode(signature),
		}, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not encode JSON: %w", err)
	}
	response, err := post("/vault/authenticated_object/"+b64encode(keyPair.publicKey), "application/json", string(jsonData))
	if err != nil {
		return fmt.Errorf("Could not write authenticated data: %s - %w", response, err)
	}
	return nil
}

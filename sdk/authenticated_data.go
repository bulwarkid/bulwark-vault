package sdk

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type AuthenticatedDataJson struct {
	Data string `json:"data"`
	Iv string `json:"iv"`
	Signature string `json:"signature"`
}

func getDataByPath(publicKey ed25519.PublicKey, encryptionKey AESEncryptionKey) ([]byte, error) {
	publicKeyEncoded := base64.URLEncoding.EncodeToString(publicKey)
	jsonData, err := get("/vault/authenticated_object/"+publicKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve data: %w",err)
	}
	dataBlob := AuthenticatedDataJson{}
	err = json.Unmarshal([]byte(jsonData), &dataBlob)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	data, err := base64.URLEncoding.DecodeString(dataBlob.Data)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Base64: %w", err)
	}
	iv, err := base64.URLEncoding.DecodeString(dataBlob.Iv)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Base64: %w", err)
	}
	signature, err := base64.URLEncoding.DecodeString(dataBlob.Signature)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Base64: %w", err)
	}
	if !verifySignature(publicKey, data, signature) {
		return nil, fmt.Errorf("Authentication failed for data for public key %s",publicKeyEncoded)
	}
	decryptedData, err := decryptBytes(data, encryptionKey, iv)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting bytes: %w", err)
	}
	return decryptedData, nil
}

func writeDataToPath(data []byte, keyPair *PublicKeyPair, encryptionKey AESEncryptionKey) error {
	signature := signData(keyPair.privateKey, data)
	encryptedData, nonce, err := encryptBytes(data, encryptionKey)
	if err != nil {
		return fmt.Errorf("Could not encrypt bytes: %w", err)
	}
	encodedData := base64.URLEncoding.EncodeToString(encryptedData)
	encodedIv := base64.URLEncoding.EncodeToString(nonce)
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	jsonData, err := json.MarshalIndent(AuthenticatedDataJson{Data: encodedData, Iv: encodedIv, Signature: encodedSignature}, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not encode JSON: %w", err)
	}
	publicKeyEncoded := base64.URLEncoding.EncodeToString(keyPair.publicKey)
	response, err := post("/vault/authenticated_object/" + publicKeyEncoded, "application/json", string(jsonData))
	if err != nil {
		return fmt.Errorf("Could not write authenticated data: %s - %w", response, err)
	}
	return nil
}

package sdk

import (
	"bytes"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	randomKeyPair, err := randomPublicKeyPair()
	checkTestError(t, err, "Could not generate public key")
	data := []byte("test")
	signature := signData(randomKeyPair.privateKey, data)
	publicKeyEncoded := b64encode(randomKeyPair.publicKey)
	publicKeyDecoded := b64decode(publicKeyEncoded)
	verified := verifySignature(publicKeyDecoded, data, signature)
	if !verified {
		t.Fatalf("Could not verify signature\n")
	}
}

func TestGetAuthenticatedData(t *testing.T) {
	randomKeyPair, err := randomPublicKeyPair()
	checkTestError(t, err, "Could not generate public key")
	encryptionKey, err := randomEncryptionKey()
	checkTestError(t, err, "Could not generate encryption key")
	data := []byte("test")
	err = writeAuthData(data, randomKeyPair, encryptionKey)
	checkTestError(t, err, "Could not write authenticated data to path")
	returnedData, err := getAuthData(randomKeyPair.publicKey, encryptionKey)
	checkTestError(t, err, "Could not retrieve auth data")
	if !bytes.Equal(data, returnedData) {
		t.Fatalf("Returned data and original data are different: %s != %s", data, returnedData)
	}
}

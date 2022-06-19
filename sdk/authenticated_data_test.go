package sdk

import (
	"bytes"
	"testing"
)

func TestGetAuthenticatedData(t *testing.T) {
	randomKeyPair,err := randomPublicKeyPair()
	checkTestError(t,err,"Could not generate public key")
	encryptionKey, err := randomEncryptionKey()
	checkTestError(t, err, "Could not generate encryption key")
	data := []byte("test")
	err = writeAuthData(data,randomKeyPair, encryptionKey)
	checkTestError(t, err, "Could not write authenticated data to path")
	returnedData, err := getAuthData(randomKeyPair.publicKey, encryptionKey)
	checkTestError(t, err, "Could not retrieve auth data")
	if bytes.Equal(data,returnedData) {
		t.Fatalf("Returned data and original data are different: %s != %s", data, returnedData)
	}
}
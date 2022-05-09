package main

import "testing"

func TestSha256(t *testing.T) {
	hash := hashSha256([]byte("test"))
	if len(hash) != 32 {
		t.Fatal("Hash wrong size")
	}
}

func TestRandomBytes(t *testing.T) {
	bytes, err := randomBytes(32)
	checkTestError(t, err, "Error generating bytes")
	if len(bytes) != 32 {
		t.Fatal("Wrong number of bytes")
	}
}

func TestBlockEncryption(t *testing.T) {
	objectData := "test_data"
	key, err := randomBytes(32)
	checkTestError(t, err, "Could not generate key")
	blob, iv, err := encryptBytes([]byte(objectData), key)
	checkTestError(t, err, "Could not encrypt data")
	outputData, err := decryptBytes(blob, key, iv)
	checkTestError(t, err, "Could not decrypt data")
	if objectData != string(outputData) {
		t.Fatalf("Output data does not match input data: (%s) -> (%s)", objectData, outputData)
	}
}

func TestLowEntropyKeygen(t *testing.T) {
	salt, err := randomBytes(32)
	checkTestError(t, err, "Could not generate random bytes")
	bytes := bytesFromLowEntropy("email:password", salt, 32)
	if len(bytes) != 32 {
		t.Fatal("Wrong number of bytes on output bytes")
	}
}

func TestHighEntropyKeygn(t *testing.T) {
	bytes, err := bytesFromHighEntropy("secret:blah:/path", 32)
	checkTestError(t, err, "Could not generate bytes")
	if len(bytes) != 32 {
		t.Fatal("Wrong number of bytes")
	}
}

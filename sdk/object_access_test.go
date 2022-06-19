package sdk

import (
	"reflect"
	"testing"
)

func TestGetSalt(t *testing.T) {
	salt, err := getSalt("email")
	checkTestError(t, err, "Failed to get salt")
	if len(salt) != 32 {
		t.Fatal("Invalid salt:", salt, err)
	}
	salt2, err := getSalt("email")
	checkTestError(t, err, "Failed to get second salt")
	if len(salt2) != 32 {
		t.Fatal("Invalid salt:", salt, err)
	}
	if !reflect.DeepEqual(salt, salt2) {
		t.Fatal("Salt is not saved between fetches:", salt, salt2)
	}
}

func TestGetObject(t *testing.T) {
	masterSecret := testMasterSecret(t)
	bytes, err := randomBytes(32)
	checkTestError(t, err, "Could not generate bytes")
	objectData := b64encode(bytes)
	err = writeObjectByPath(masterSecret, "/test", objectData)
	checkTestError(t, err, "Could not write object")
	returnedData, err := getObjectByPath(masterSecret, "/test")
	checkTestError(t, err, "Could not retrive object")
	if objectData != returnedData {
		t.Logf("Returned data is different than written data: (%s) -> (%s)", objectData, returnedData)
	}
}

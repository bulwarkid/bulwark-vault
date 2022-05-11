package sdk

import "testing"

func testMasterSecret(t *testing.T) []byte {
	loginSecret, err := deriveLoginSecret("email", "password")
	checkTestError(t, err, "Could not derive login secret")
	masterSecret, err := getMasterSecret(loginSecret)
	checkTestError(t, err, "Could not get master secret")
	return masterSecret
}

func checkTestError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %s", msg, err)
	}
}

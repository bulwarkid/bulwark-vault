package sdk

import (
	"fmt"
)

const MASTER_SECRET_SIZE = 32
const LOGIN_SECRET_SIZE = 32

func deriveLoginSecret(email string, password string) ([]byte, error) {
	salt, err := getSalt(email)
	if err != nil {
		return nil, err
	}
	return bytesFromLowEntropy("login_secret:"+email+":"+password, salt, LOGIN_SECRET_SIZE), nil
}

func getMasterSecret(loginSecret []byte) ([]byte, error) {
	secretBase64, err := getObjectByPath(loginSecret, "/master-secret")
	if err != nil {
		if !isReturnCode(err, 404) {
			return nil, fmt.Errorf("Could not get master secret: %w", err)
		}
		secretBytes, err := newMasterSecret()
		if err != nil {
			return nil, fmt.Errorf("Could not generate random bytes: %w", err)
		}
		secretBase64 = b64encode(secretBytes)
		if err = writeObjectByPath(loginSecret, "/master-secret", b64encode(secretBytes)); err != nil {
			return nil, fmt.Errorf("Could not write master secret: %w", err)
		}
		return secretBytes, nil
	}
	return b64decode(secretBase64), nil
}

func newMasterSecret() ([]byte, error) {
	return randomBytes(MASTER_SECRET_SIZE)
}

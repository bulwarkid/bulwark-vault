package sdk

import (
	"fmt"
)

type VaultAccess struct {
	masterSecret []byte
	keyDirectory *KeyDirectory
}

func NewVault() *VaultAccess {
	defer recoverFromPanic()
	return &VaultAccess{masterSecret: nil, keyDirectory: nil}
}

func (access *VaultAccess) Login(email, password string) error {
	defer recoverFromPanic()
	loginSecret, err := deriveLoginSecret(email, password)
	if err != nil {
		return err
	}
	masterSecret, err := getMasterSecret(loginSecret)
	if err != nil {
		return err
	}
	keyDirectory, err := getKeyDirectory(access.masterSecret)
	if err != nil {
		return err
	}
	access.masterSecret = masterSecret
	access.keyDirectory = keyDirectory
	return nil
}

func (access *VaultAccess) Put(path string, data string) error {
	defer recoverFromPanic()
	var err error
	if access.keyDirectory == nil {
		return fmt.Errorf("Vault isn't logged in")
	}
	accessData, err := access.keyDirectory.getOrCreateForPath(path)
	if err != nil {
		return fmt.Errorf("Could not get or create for path %s: %w", path, err)
	}
	if err = writeObject(accessData.accessKey, accessData.encryptionKey, data); err != nil {
		return fmt.Errorf("Could not write object: %w", err)
	}
	return nil
}

func (access *VaultAccess) Get(path string) (string, error) {
	defer recoverFromPanic()
	if access.keyDirectory == nil {
		return "", fmt.Errorf("Vault isn't logged in")
	}
	accessData, err := access.keyDirectory.getPath(path)
	if err != nil {
		return "", fmt.Errorf("No value at path: %s", path)
	}
	data, err := getObject(accessData.accessKey, accessData.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("Could not get object: %w", err)
	}
	return data, nil
}

func (vault *VaultAccess) MasterSecret() string {
	defer recoverFromPanic()
	if vault.masterSecret == nil {
		return ""
	}
	return b64encode(vault.masterSecret)
}

func (vault *VaultAccess) KeyDirectory() *KeyDirectory {
	defer recoverFromPanic()
	return vault.keyDirectory
}

func GetAuthData(publicKeyBase64 string, encryptionKeyBase64 string) ([]byte,error) {
	publicKey := b64decode(publicKeyBase64)
	encryptionKey := b64decode(encryptionKeyBase64)
	bytes, err := getAuthData(publicKey,encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Auth Data: %w",err)
	}
	return bytes, nil
}

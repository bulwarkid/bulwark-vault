package sdk

import (
	"fmt"
)

type VaultAccess struct {
	masterSecret []byte
	keyDirectory *KeyDirectory
}

func NewVault() *VaultAccess {
	return &VaultAccess{masterSecret: nil, keyDirectory: nil}
}

func (access *VaultAccess) login(email, password string) error {
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

func (access *VaultAccess) store(path string, data string) error {
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

func (access *VaultAccess) get(path string) (string, error) {
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

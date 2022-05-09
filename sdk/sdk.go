package main

type VaultAccess struct {
	masterSecret []byte
	keyDirectory *KeyDirectory
}

func newVault() *VaultAccess {
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
	access.masterSecret = masterSecret
	keyDirectory, err := getKeyDirectory(access.masterSecret)
	if err != nil {
		return err
	}
	access.keyDirectory = keyDirectory
	return nil
}

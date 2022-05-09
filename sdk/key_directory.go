package main

import "encoding/json"

type AccessData struct {
	accessKey     string
	encryptionKey string
}

type KeyDirectory map[string]*AccessData

func NewKeyDirectory() *KeyDirectory {
	directory := make(KeyDirectory)
	return &directory
}

func (directory *KeyDirectory) Load(masterSecret []byte) error {
	keyDirectoryObject, err := getObjectByPath(masterSecret, "/directory")
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(keyDirectoryObject), directory); err != nil {
		return err
	}
	return nil
}

func (directory *KeyDirectory) Store(masterSecret []byte) error {
	jsonData, err := json.Marshal(directory)
	if err != nil {
		return err
	}
	return writeObjectByPath([]byte(jsonData), "/master-secret", string(jsonData))
}

func getKeyDirectory(masterSecret []byte) (*KeyDirectory, error) {
	directory := NewKeyDirectory()
	err := directory.Load(masterSecret)
	if err != nil {
		if !isReturnCode(err, 404) {
			return nil, err
		}
		if err = directory.Store(masterSecret); err != nil {
			return nil, err
		}
	}
	return directory, nil
}

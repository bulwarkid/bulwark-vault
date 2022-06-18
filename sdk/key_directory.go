package sdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type AccessData struct {
	accessKey     AccessKey
	encryptionKey AESEncryptionKey
}

func randomAccessData() (*AccessData, error) {
	accessKey, err := randomAccessKey()
	if err != nil {
		return nil, fmt.Errorf("Could not generate access key: %w", err)
	}
	encryptionKey, err := randomEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("Could not generate encryption key: %w", err)
	}
	return &AccessData{accessKey, encryptionKey}, nil
}

type KeyDirectory struct {
	masterSecret []byte
	values       *map[string]*AccessData
}

func newKeyDirectory(masterSecret []byte) *KeyDirectory {
	directory := KeyDirectory{}
	directory.masterSecret = masterSecret
	values := make(map[string]*AccessData)
	directory.values = &values
	return &directory
}

type AccessDataJson struct {
	AccessKey     string `json:"accessKey"`
	EncryptionKey string `json:"encryptionKey"`
}

func (directory *KeyDirectory) Json() (string, error) {
	var keyDirectoryJson map[string]AccessDataJson = make(map[string]AccessDataJson)
	for path, accessData := range *directory.values {
		encryptionKeyEncoded := base64.URLEncoding.EncodeToString(accessData.encryptionKey)
		keyDirectoryJson[path] = AccessDataJson{AccessKey: string(accessData.accessKey), EncryptionKey: encryptionKeyEncoded}
	}
	jsonData, err := json.MarshalIndent(keyDirectoryJson, "", "    ")
	if err != nil {
		return "", fmt.Errorf("Could not encode JSON: %w", err)
	}
	return string(jsonData), nil
}

func (directory *KeyDirectory) load() error {
	keyDirectoryObject, err := getObjectByPath(directory.masterSecret, "/directory")
	if err != nil {
		return err
	}
	var keyDirectoryJson map[string]AccessDataJson
	if err := json.Unmarshal([]byte(keyDirectoryObject), &keyDirectoryJson); err != nil {
		return err
	}
	values := make(map[string]*AccessData)
	for path, accessDataJson := range keyDirectoryJson {
		encryptionKeyDecoded, err := base64.URLEncoding.DecodeString(accessDataJson.EncryptionKey)
		if err != nil {
			return fmt.Errorf("Error decoding encryption keys: %w", err)
		}
		values[path] = &AccessData{accessKey: AccessKey(accessDataJson.AccessKey), encryptionKey: AESEncryptionKey(encryptionKeyDecoded)}
	}
	directory.values = &values
	return nil
}

func (directory *KeyDirectory) store() error {
	jsonData, err := directory.Json()
	if err != nil {
		return fmt.Errorf("Could not encode JSON: %w", err)
	}
	if err = writeObjectByPath(directory.masterSecret, "/directory", string(jsonData)); err != nil {
		return fmt.Errorf("Could not write key directory: %w", err)
	}
	return nil
}

func (directory *KeyDirectory) getPath(path string) (*AccessData, error) {
	accessData, ok := (*directory.values)[path]
	if !ok {
		return nil, fmt.Errorf("No value for path: %s", path)
	}
	return accessData, nil
}

func (directory *KeyDirectory) getOrCreateForPath(path string) (*AccessData, error) {
	accessData, err := directory.getPath(path)
	if err == nil {
		return accessData, nil
	}
	accessData, err = randomAccessData()
	if err != nil {
		return nil, fmt.Errorf("Could not generate random access data: %w", err)
	}
	if err = directory.writePath(path, *accessData); err != nil {
		return nil, fmt.Errorf("Could not store generated data: %w", err)
	}
	return accessData, nil
}

func (directory *KeyDirectory) writePath(path string, accessData AccessData) error {
	(*directory.values)[path] = &accessData
	if err := directory.store(); err != nil {
		delete(*directory.values, path)
		return fmt.Errorf("Could not store generated data: %w", err)
	}
	return nil
}

func getKeyDirectory(masterSecret []byte) (*KeyDirectory, error) {
	directory := newKeyDirectory(masterSecret)
	err := directory.load()
	if err != nil {
		if !isReturnCode(err, 404) {
			return nil, err
		}
		if err = directory.store(); err != nil {
			return nil, err
		}
	}
	return directory, nil
}

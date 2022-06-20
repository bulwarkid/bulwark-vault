package sdk

import (
	"encoding/json"
	"fmt"
)

const ACCESS_KEY_LENGTH = 32

type AccessKey string

func randomAccessKey() (AccessKey, error) {
	bytes, err := randomBytes(ACCESS_KEY_LENGTH)
	if err != nil {
		return "", fmt.Errorf("Could not generate access key: %w", err)
	}
	return AccessKey(b64encode(bytes)), nil
}

func saltId(data string) string {
	saltIdBytes := hashSha256([]byte(data))
	return b64encode(saltIdBytes)
}

func getSalt(inputData string) ([]byte, error) {
	saltId := saltId(inputData)
	saltBase64, err := get("/vault/salt/" + saltId)
	if err != nil {
		return nil, err
	}
	return b64decode(saltBase64), nil
}

type EncryptedBlob struct {
	Blob string `json:"blob"`
	Iv   string `json:"iv"`
}

func encodeBlob(blob []byte, iv []byte) (string, error) {
	blobJson, err := json.Marshal(
		&EncryptedBlob{
			Blob: b64encode(blob),
			Iv:   b64encode(iv),
		})
	if err != nil {
		return "", err
	}
	return string(blobJson), nil
}

func decodeBlob(blobJson string) ([]byte, []byte, error) {
	var blob EncryptedBlob
	if err := json.Unmarshal([]byte(blobJson), &blob); err != nil {
		return nil, nil, err
	}
	return b64decode(blob.Blob), b64decode(blob.Iv), nil
}

func deterministicObjectAccessKey(secret []byte, path string) (AccessKey, error) {
	bytes, err := bytesFromHighEntropy("object_access_key:"+b64encode(secret)+":"+path, ACCESS_KEY_LENGTH)
	if err != nil {
		return "", fmt.Errorf("Could not generate access key: %w", err)
	}
	return AccessKey(b64encode(bytes)), nil
}

func getObjectByPath(secret []byte, path string) (string, error) {
	accessKey, err := deterministicObjectAccessKey(secret, path)
	if err != nil {
		return "", fmt.Errorf("Could not get access key: %w", err)
	}
	encryptionKey, err := deterministicObjectEncryptionKey(secret, path)
	if err != nil {
		return "", fmt.Errorf("Could not get encryption key: %w", err)
	}
	return getObject(accessKey, encryptionKey)
}

func getObject(accessKey AccessKey, encryptionKey AESEncryptionKey) (string, error) {
	objectBase64, err := get("/vault/object/" + string(accessKey))
	if err != nil {
		return "", fmt.Errorf("Could not retrieve object data: %w", err)
	}
	blobJson := b64decode(objectBase64)
	blob, iv, err := decodeBlob(string(blobJson))
	if err != nil {
		return "", fmt.Errorf("Could not decode blob json: %w (%s)", err, blobJson)
	}
	decryptedBytes, err := decryptBytes(blob, encryptionKey, iv)
	if err != nil {
		return "", fmt.Errorf("Could not decrypt bytes: %w", err)
	}
	return string(decryptedBytes), nil
}

func writeObjectByPath(secret []byte, path string, data string) error {
	accessKey, err := deterministicObjectAccessKey(secret, path)
	if err != nil {
		return err
	}
	encryptionKey, err := deterministicObjectEncryptionKey(secret, path)
	if err != nil {
		return err
	}
	return writeObject(accessKey, encryptionKey, data)
}

func writeObject(accessKey AccessKey, encryptionKey AESEncryptionKey, data string) error {
	encryptedBytes, iv, err := encryptBytes([]byte(data), encryptionKey)
	if err != nil {
		return fmt.Errorf("Could not encrypt data: %w", err)
	}
	blobJson, err := encodeBlob(encryptedBytes, iv)
	if err != nil {
		return fmt.Errorf("Could not create json: %w", err)
	}
	objectBase64 := b64encode([]byte(blobJson))
	_, err = post("/vault/object/"+string(accessKey), "text/plain", objectBase64)
	if err != nil {
		return fmt.Errorf("Could not store encrypted data: %w", err)
	}
	return nil
}

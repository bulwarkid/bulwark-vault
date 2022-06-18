package sdk

import (
	"encoding/base64"
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
	return AccessKey(base64.URLEncoding.EncodeToString(bytes)), nil
}

func saltId(data string) string {
	saltIdBytes := hashSha256([]byte(data))
	return base64.URLEncoding.EncodeToString(saltIdBytes)
}

func getSalt(inputData string) ([]byte, error) {
	saltId := saltId(inputData)
	saltBase64, err := get("/vault/salt/" + saltId)
	if err != nil {
		return nil, err
	}
	salt, err := base64.URLEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

type EncryptedBlob struct {
	Blob string `json:"blob"`
	Iv   string `json:"iv"`
}

func encodeBlob(blob []byte, iv []byte) (string, error) {
	blobBase64 := base64.URLEncoding.EncodeToString(blob)
	ivBase64 := base64.URLEncoding.EncodeToString(iv)
	blobJson, err := json.Marshal(&EncryptedBlob{Blob: blobBase64, Iv: ivBase64})
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
	blobBytes, err := base64.URLEncoding.DecodeString(blob.Blob)
	if err != nil {
		return nil, nil, err
	}
	ivBytes, err := base64.URLEncoding.DecodeString(blob.Iv)
	if err != nil {
		return nil, nil, err
	}
	return blobBytes, ivBytes, nil
}

func deterministicObjectAccessKey(secret []byte, path string) (AccessKey, error) {
	secretString := base64.URLEncoding.EncodeToString(secret)
	bytes, err := bytesFromHighEntropy("object_access_key:"+secretString+":"+path, ACCESS_KEY_LENGTH)
	if err != nil {
		return "", fmt.Errorf("Could not generate access key: %w", err)
	}
	return AccessKey(base64.URLEncoding.EncodeToString(bytes)), nil
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
	blobJson, err := base64.URLEncoding.DecodeString(objectBase64)
	if err != nil {
		return "", fmt.Errorf("Could not decode object bytes: %w (%s)", err, objectBase64)
	}
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
	objectBase64 := base64.URLEncoding.EncodeToString([]byte(blobJson))
	_, err = post("/vault/object/"+string(accessKey), "text/plain", objectBase64)
	if err != nil {
		return fmt.Errorf("Could not store encrypted data: %w", err)
	}
	return nil
}

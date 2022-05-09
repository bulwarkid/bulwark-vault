package main

import "testing"

func TestGetKeyDirectory(t *testing.T) {
	masterSecret := testMasterSecret(t)
	directory, err := getKeyDirectory(masterSecret)
	checkTestError(t, err, "Error getting key directory")
	var data AccessData
	data.accessKey = "access"
	data.encryptionKey = "encryption"
	(*directory)["/test"] = &data
	err = directory.Store(masterSecret)
	checkTestError(t, err, "Error storing directory")
}

package main

import "testing"

func TestGetKeyDirectory(t *testing.T) {
	masterSecret := testMasterSecret(t)
	directory, err := getKeyDirectory(masterSecret)
	checkTestError(t, err, "Error getting key directory")
	data, err := randomAccessData()
	checkTestError(t, err, "Error generating random access data")
	err = directory.writePath("/test", *data)
	checkTestError(t, err, "Error storing directory")
}

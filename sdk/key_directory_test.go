package sdk

import "testing"

func TestGetKeyDirectory(t *testing.T) {
	masterSecret := testMasterSecret(t)
	directory, err := getKeyDirectory(masterSecret)
	checkTestError(t, err, "Error getting key directory")
	data, err := randomAccessData()
	checkTestError(t, err, "Error generating random access data")
	err = directory.writePath("/test", *data)
	checkTestError(t, err, "Error storing directory")
	outputData, err := directory.getPath("/test")
	checkTestError(t, err, "Error getting data")
	if data.accessKey != outputData.accessKey {
		t.Fatal("Output data is incorrect")
	}
	directory2, err := getKeyDirectory(masterSecret)
	checkTestError(t, err, "Error getting directory again")
	outputData2, err := directory2.getPath("/test")
	checkTestError(t, err, "Error getting data from second directory")
	if outputData.accessKey != outputData2.accessKey {
		t.Fatal("Output data does not match from two directories")
	}
}

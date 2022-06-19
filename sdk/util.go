package sdk

import (
	"encoding/base64"
	"fmt"
	"runtime/debug"
)

func b64encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func b64decode(dataString string) []byte {
	data, err := base64.URLEncoding.DecodeString(dataString)
	if err != nil {
		panic(fmt.Sprintf("Could not decode base64: %s - %s", dataString, debug.Stack()))
	}
	return data
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic:", r)
	}
}

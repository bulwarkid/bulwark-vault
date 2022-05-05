package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func generateHash() string {
	hash := make([]byte, 32)
	rand.Read(hash)
	return base64.URLEncoding.EncodeToString(hash)
}

func verifySalt(salt string) bool {
	bytes, err := base64.URLEncoding.DecodeString(salt)
	return err == nil && len(bytes) == 32
}

func testSalt(iterations int) {
	hashes := make([]string, 10000)
	for i := range hashes {
		hashes[i] = generateHash()
	}
	start := float64(time.Now().UnixMilli())
	totalRuns := 0
	for i := 0; i < iterations; i++ {
		for _, hash := range hashes {
			response, err := http.Get("http://localhost:8080/vault/salt/" + hash)
			if err != nil {
				fmt.Println("Bad server response:", err)
				return
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil || !verifySalt(string(body)) {
				fmt.Println("Bad salt returned:", string(body))
				return
			}
		}
		totalRuns += len(hashes)
		current := float64(time.Now().UnixMilli())
		seconds := (current - start) / 1000.0
		fmt.Println("Requests/sec:", float64(totalRuns)/seconds)
	}
}

func main() {
	testSalt(10)
}

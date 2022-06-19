package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func handleSalt(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	saltId := pathParts[len(pathParts)-1]
	if !verifyHash(saltId) {
		returnCode(w, 403)
		return
	}
	salt, err := getSalt(saltId)
	if err != nil {
		salt, err = generateSalt(saltId)
		if err != nil {
			fmt.Println("Error: ", err)
			returnCode(w, 500)
			return
		}
	}
	fmt.Fprintf(w, salt)
}

func verifyHash(hash string) bool {
	bytes, err := base64.URLEncoding.DecodeString(hash)
	return err == nil && len(bytes) == 32
}

func getSalt(saltId string) (string, error) {
	rows, err := db.Query(`SELECT salt FROM salts WHERE salt_id=$1`, saltId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var salt string
		err = rows.Scan(&salt)
		if err != nil {
			return "", err
		}
		return salt, nil
	}
	return "", fmt.Errorf("No salt found")
}

func generateSalt(saltId string) (string, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", fmt.Errorf("Could not generate random bytes for salt")
	}
	salt := base64.URLEncoding.EncodeToString(saltBytes)
	_, err = db.Exec(`INSERT INTO salts(salt_id,salt) values($1,$2)`, saltId, salt)
	if err != nil {
		return "", err
	}
	return salt, nil
}

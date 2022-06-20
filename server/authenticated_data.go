package main

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthDataJson struct {
	Data string `json:"data"`
	Iv string `json:"iv"`
	Signature string `json:"signature"`
}

func verifyData(jsonBlob string, publicKeyBase64 string) (bool, error) {
	jsonData := AuthDataJson{}
	err := json.Unmarshal([]byte(jsonBlob), &jsonData)
	if err != nil {
		return false, fmt.Errorf("Error decoding Auth Data JSON: %w - %s", err, jsonBlob)
	}
	publicKey, err := base64.URLEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, fmt.Errorf("Could not decode public key: %w", err)
	}
	data, err := base64.URLEncoding.DecodeString(jsonData.Data)
	if err != nil {
		return false, fmt.Errorf("Invalid Base64 for data: %w - %s", err, jsonData.Data)
	}
	signature, err := base64.URLEncoding.DecodeString(jsonData.Signature)
	if err != nil {
		return false, fmt.Errorf("Invalid Base64 for signature: %w - %s", err, jsonData.Signature)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return false, fmt.Errorf("Wrong public key size when verifying data")
	}
	return ed25519.Verify(publicKey, data, signature), nil
}

func getAuthData(publicKeyBase64 string) (string, error) {
	db := getDb()
	rows, err := db.Query("SELECT object_data FROM authenticated_objects WHERE object_id=$1", publicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("Could not read authenticated data: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var jsonBlob string
		err = rows.Scan(&jsonBlob)
		if err != nil {
			return "", fmt.Errorf("Could not read values in row: %w",err)
		}
		return jsonBlob, nil
	}
	return "", nil
}

func writeAuthData(publicKeyBase64 string, jsonBlob string) error {
	verified, err := verifyData(jsonBlob,publicKeyBase64)
	if err != nil || !verified {
		return fmt.Errorf("Could not verify signature on authenticated data")
	}
	ctx := context.Background()
	tx, err := getDb().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Could not write authenticated data: %w", err)
	}
	row := tx.QueryRowContext(ctx, `SELECT object_data FROM authenticated_objects WHERE object_id=$1`, publicKeyBase64)
	var oldData string
	err = row.Scan(&oldData)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return fmt.Errorf("Could not write authenticated data: %w", err)
	}
	if err != sql.ErrNoRows {
		_, err := tx.ExecContext(ctx, `UPDATE authenticated_objects SET object_data=$1 WHERE object_id=$2`, jsonBlob, publicKeyBase64)
		if err != nil {
			fmt.Println("Rolling back:", err)
			tx.Rollback()
			return fmt.Errorf("Could not write authenticated data: %w", err)
		}
	} else {
		_, err := tx.Exec(`INSERT INTO authenticated_objects(object_id,object_data) values($1,$2)`, publicKeyBase64, jsonBlob)
		if err != nil {
			fmt.Println("Rolling back:", err)
			tx.Rollback()
			return fmt.Errorf("Could not write authenticated data: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Could not write authenticated data: %w", err)
	}
	return nil
}

func handleAuthDataGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	publicKeyBase64 := pathParts[len(pathParts)-1]
	jsonBlob, err := getAuthData(publicKeyBase64)
	if err != nil {
		fmt.Printf("Server error: %s\n",err)
		returnCode(w, 500)
		return
	}
	if jsonBlob == "" {
		returnCode(w, 404)
		return
	}
	verified, err := verifyData(jsonBlob,publicKeyBase64)
	if err != nil || !verified {
		returnCode(w, 403)
		return
	}
	fmt.Fprintf(w, jsonBlob)
}

func handleAuthDataPost(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	publicKeyBase64 := pathParts[len(pathParts)-1]
	jsonBlob, err := ioutil.ReadAll(io.LimitReader(r.Body, 10000))
	if err != nil {
		fmt.Println("Error reading user data")
		returnCode(w, 400)
		return
	}
	verified, err := verifyData(string(jsonBlob),publicKeyBase64)
	if err != nil || !verified {
		fmt.Printf("Could not verify authenticated data: %s %v\n", err,verified)
		returnCode(w, 403)
		return
	}
	err = writeAuthData(publicKeyBase64, string(jsonBlob))
	if err != nil {
		fmt.Printf("Server Error: %s\n", err)
		returnCode(w, 500)
		return
	}
	fmt.Fprintf(w, "Success")
}
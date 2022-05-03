package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "vault"
)

var db *sql.DB

func loadDb() (*sql.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
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

func getObject(objectId string) (string, error) {
	rows, err := db.Query(`SELECT object_data FROM objects WHERE object_id=$1`, objectId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var data string
		err = rows.Scan(&data)
		if err != nil {
			return "", err
		}
		return data, nil
	}
	return "", fmt.Errorf("No salt found")
}

func returnCode(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func handleSalt(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	saltId := pathParts[len(pathParts)-1]
	if len(saltId) != 44 {
		returnCode(w, 403, "Invalid Salt")
	}
	salt, err := getSalt(saltId)
	if err != nil {
		salt, err = generateSalt(saltId)
		if err != nil {
			returnCode(w, 500, "Server error")
		}
	}
	fmt.Fprintf(w, salt)
}

func handleObject(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	objectId := pathParts[len(pathParts)-1]
	if len(objectId) != 44 {
		// TODO: Return 400
		returnCode(w, 403, "Invalid Object ID")
	}
	objectData, err := getObject(objectId)
	if err != nil {
		returnCode(w, 404, "404 - Object not found")
	}
	fmt.Fprintf(w, objectData)
}

func main() {
	fmt.Println("Starting vault server...")
	var err error
	db, err = loadDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to DB!")

	http.HandleFunc("/vault/salt/", handleSalt)
	http.HandleFunc("/vault/object/", handleObject)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

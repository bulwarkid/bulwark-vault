package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
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
	return "", fmt.Errorf("No object found")
}

func writeObject(objectId string, objectData string) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	row := tx.QueryRowContext(ctx, `SELECT object_data FROM objects WHERE object_id=$1`, objectId)
	var oldData string
	err = row.Scan(&oldData)
	if err != nil && err != sql.ErrNoRows{
		tx.Rollback()
		return err
	}
	if err != sql.ErrNoRows {
		_, err := tx.ExecContext(ctx,`UPDATE objects SET object_data=$1 WHERE object_id=$2`, objectData, objectId)
		if err != nil {
			fmt.Println("Rolling back:",err)
			tx.Rollback()
			return err
		}
	} else {
		_, err := tx.Exec(`INSERT INTO objects(object_id,object_data) values($1,$2)`, objectId, objectData)
		if err != nil {
			fmt.Println("Rolling back:",err)
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func returnCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%d - %s", code, http.StatusText(code))))
}

type httpHandler func(w http.ResponseWriter, r *http.Request)

func requestHandler(getHandler httpHandler, postHandler httpHandler) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin","*")
		switch r.Method {
		case "GET":
			if getHandler == nil {
				returnCode(w, 405)
			} else {
				getHandler(w, r)
			}
		case "POST":
			if postHandler == nil {
				returnCode(w, 405)
			} else {
				postHandler(w, r)
			}
		default:
			returnCode(w, 405)
		}
	}
}

func readLimit(r io.Reader, limit int64) ([]byte, error) {
	output, err := ioutil.ReadAll(io.LimitReader(r, limit))
	if err != nil {
		return nil, err
	}
	return output, nil
}

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
			returnCode(w, 500)
			return
		}
	}
	fmt.Fprintf(w, salt)
}

func handleObjectGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	objectId := pathParts[len(pathParts)-1]
	if !verifyHash(objectId) {
		returnCode(w, 403)
		return
	}
	objectData, err := getObject(objectId)
	if err != nil {
		returnCode(w, 404)
		return
	}
	fmt.Fprintf(w, objectData)
}

func handleObjectPost(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	objectId := pathParts[len(pathParts)-1]
	if !verifyHash(objectId) {
		returnCode(w, 403)
	}
	objectBytes, err := ioutil.ReadAll(io.LimitReader(r.Body, 10000))
	if err != nil {
		returnCode(w, 404)
		return
	}
	objectData := string(objectBytes)
	writeObject(objectId, objectData)
	fmt.Fprintf(w, "Success")
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

	http.HandleFunc("/vault/salt/", requestHandler(handleSalt, nil))
	http.HandleFunc("/vault/object/", requestHandler(handleObjectGet, handleObjectPost))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

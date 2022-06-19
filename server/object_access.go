package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func getObject(objectId string) (string, error) {
	rows, err := getDb().Query(`SELECT object_data FROM objects WHERE object_id=$1`, objectId)
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
	tx, err := getDb().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	row := tx.QueryRowContext(ctx, `SELECT object_data FROM objects WHERE object_id=$1`, objectId)
	var oldData string
	err = row.Scan(&oldData)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}
	if err != sql.ErrNoRows {
		_, err := tx.ExecContext(ctx, `UPDATE objects SET object_data=$1 WHERE object_id=$2`, objectData, objectId)
		if err != nil {
			fmt.Println("Rolling back:", err)
			tx.Rollback()
			return err
		}
	} else {
		_, err := tx.Exec(`INSERT INTO objects(object_id,object_data) values($1,$2)`, objectId, objectData)
		if err != nil {
			fmt.Println("Rolling back:", err)
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

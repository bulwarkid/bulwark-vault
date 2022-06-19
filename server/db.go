package main

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "vault_user"
	password = "insecure_password"
	dbname   = "vault"
)

var db *sql.DB

func loadDb() error {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		db = nil
		return err
	}
	err = db.Ping()
	if err != nil {
		db = nil
		return err
	}
	return nil
}

func closeDb() {
	if db != nil {
		db.Close()
	}
}

func getDb() *sql.DB {
	if db == nil {
		panic("Trying to access DB before connection established!")
	}
	return db
}

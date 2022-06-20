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

var internalDb *sql.DB

func loadDb() error {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	internalDb, err = sql.Open("postgres", psqlConn)
	if err != nil {
		internalDb = nil
		return err
	}
	err = internalDb.Ping()
	if err != nil {
		internalDb = nil
		return err
	}
	return nil
}

func closeDb() {
	if internalDb != nil {
		internalDb.Close()
	}
}

func getDb() *sql.DB {
	if internalDb == nil {
		panic("Trying to access DB before connection established!")
	}
	return internalDb
}

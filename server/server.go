package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting vault server...")
	err := loadDb()
	if err != nil {
		panic(err)
	}
	defer closeDb()
	fmt.Println("Connected to DB!")

	http.HandleFunc("/vault/salt/", requestHandler(handleSalt, nil))
	http.HandleFunc("/vault/object/", requestHandler(handleObjectGet, handleObjectPost))
	log.Fatal(http.ListenAndServe(":5001", nil))
}

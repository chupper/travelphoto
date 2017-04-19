package controllers

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

//Page is the base ttype for the views
type Page struct {
	Title string
	Body  []byte
}

// DbConnection get the connection
func DbConnection() (*sql.DB, error) {

	// Opens the connection
	con, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Connect fail: ", err)
	}

	return con, err
}

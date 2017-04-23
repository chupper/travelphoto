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

var db *sql.DB

// DbConnection get the connection
func DbConnection() (*sql.DB, error) {

	// If already exists we are all good
	if db != nil {
		return db, nil
	}

	// Opens the connection
	log.Println("Creating new DB Connection")
	con, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Connect fail: ", err)
	}
	db = con
	return con, err

}

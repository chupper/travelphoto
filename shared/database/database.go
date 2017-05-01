package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgres
)

var db *sql.DB

// Connection is an interface for making the queries
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// DbConfig contains the database connection info
type DbConfig struct {
	Provider string
	User     string
	DbName   string
	Password string
	SSLMode  string
}

// initialises the database
func Configure(c DbConfig) {
	// Opens the connection
	log.Println("Creating new DB Connection")
	con, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Connect fail: ", err)
	}
	db = con
}

// DbConnection get the connection
func DbConnection() *sql.DB {

	// if doesn't exist we open the new one
	if db == nil {
		log.Fatal("Connect fail: ")
	}

	return db
}

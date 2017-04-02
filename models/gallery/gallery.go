package gallery

import (
	"database/sql"
	"fmt"
)

const (
	table = "gallery"
)

// Gallery model
type Gallery struct {
	ID          int
	Name        string
	Description string
}

// Connection is an interface for making the queries
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Create adds a gallery
func Create(db Connection, name string, description string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v (name, description)
		VALUES ($1, $2)
		`, table), name, description)

	return result, err
}

// Select returns all the galleries
func Select(db Connection) (*[]Gallery, error) {

	result, err := db.Query(fmt.Sprintf(`
		SELECT ID, NAME, DESCRIPTION
		FROM %v
	`, table))

	if err != nil {
		return nil, err
	}

	// not sure if this is a golang way?
	var items []Gallery
	for result.Next() {
		var id int
		var name, description string
		err = result.Scan(&id, &name, &description)

		gallery := Gallery{
			ID:          id,
			Name:        name,
			Description: description,
		}

		items = append(items, gallery)
	}

	return &items, err
}

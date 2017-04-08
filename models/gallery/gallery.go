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
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Create adds a gallery
func Create(db Connection, name string, description string) int {
	var id int
	result := db.QueryRow(fmt.Sprintf(`
		INSERT INTO %v (name, description)
		VALUES ($1, $2)
		RETURNING ID
		`, table), name, description)
	result.Scan(&id)

	return id
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

// Get returns single gallery
func Get(db Connection, galleryID int) (*Gallery, error) {

	result, err := db.Query(fmt.Sprintf(`
		SELECT ID, NAME, DESCRIPTION
		FROM %v
		WHERE ID = $1
	`, table), galleryID)

	for result.Next() {
		var id int
		var name, description string
		err = result.Scan(&id, &name, &description)

		if err == nil {
			return &Gallery{
				ID:          id,
				Name:        name,
				Description: description,
			}, nil
		}
	}

	return nil, err
}

// Update updates the table
func Update(db Connection, gallery Gallery) error {

	_, err := db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			DESCRIPTION = $1
		WHERE
			ID = $2
	`, table), gallery.Description, gallery.ID)

	return err
}

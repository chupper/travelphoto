package photo

import (
	"database/sql"
	"fmt"
	"log"
)

// Photo image in the gallery
type Photo struct {
	ID          int
	Name        string
	FileName    string
	Description string
	Rank        int
	FullSize    []byte
	Thumbnail   []byte
}

const (
	table = "photo"
)

// Connection is an interface for making the queries
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Create adds a new photo for the gallery
func Create(db Connection, name string, description string, galleryID int) int {
	var id int
	result := db.QueryRow(fmt.Sprintf(`
		INSERT INTO %v (name, description, galleryId)
		VALUES ($1, $2, $3)
		RETURNING ID
		`, table), name, description, galleryID)
	result.Scan(&id)

	return id
}

// GetPhotos get photos for specific gallery
func GetPhotos(db Connection, galleryID int) (*[]Photo, error) {

	results, err := db.Query(fmt.Sprintf(`
		SELECT
			ID,
			NAME,
			DESCRIPTION
		FROM %v
		WHERE GALLERYID = $1
		`, table), galleryID)

	if err != nil {
		log.Println(err)
	}

	// not sure if this is a golang way?
	var items []Photo
	for results.Next() {
		var id int
		var name, description string
		err = results.Scan(&id, &name, &description)

		gallery := Photo{
			ID:          id,
			Name:        name,
			Description: description,
		}

		items = append(items, gallery)
	}

	return &items, err
}

// UpdatePhoto updates the main photo
func UpdatePhoto(db Connection, photoId int, photoBytes *[]byte) {

	db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			IMAGE = $1
		WHERE ID = $2
		`, table), photoBytes, photoId)
}

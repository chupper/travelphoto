package photo

import (
	"database/sql"
	"fmt"
	"log"
)

// Photo image in the gallery
type Photo struct {
	ID            int
	GalleryID     int
	Name          string
	FileName      string
	ThumbFileName string
	Description   string
	Rank          int
	FullSize      []byte
	Thumbnail     []byte
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
		INSERT INTO %v (name, description, galleryId, thumbfilename, filename)
		VALUES ($1, $2, $3, '', '')
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
			GALLERYID,
			NAME,
			DESCRIPTION,
			FILENAME,
			THUMBFILENAME
		FROM %v
		WHERE GALLERYID = $1
		ORDER BY ID
		`, table), galleryID)

	if err != nil {
		log.Println(err)
	}

	// not sure if this is a golang way?
	var items []Photo
	for results.Next() {
		var id, galleryID int
		var name, description, fileName, thumbfilename string
		err = results.Scan(&id, &galleryID, &name, &description, &fileName, &thumbfilename)

		log.Println(name, "- thumb: ", thumbfilename, "file: ", fileName)
		gallery := Photo{
			ID:            id,
			GalleryID:     galleryID,
			Name:          name,
			FileName:      fileName,
			ThumbFileName: thumbfilename,
			Description:   description,
		}

		items = append(items, gallery)
	}

	return &items, err
}

// UpdatePhoto updates the main photo
func UpdatePhoto(db Connection, photoID int, fileName string, photoBytes *[]byte) {
	db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			IMAGE = $1,
			FILENAME = $2
		WHERE ID = $3
		`, table), photoBytes, fileName, photoID)
}

// UpdateThumb updates the thumbnail
func UpdateThumb(db Connection, photoID int, fileName string, photoBytes *[]byte) {
	db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			Thumb = $1,
			ThumbFileName = $2
		WHERE ID = $3
		`, table), photoBytes, fileName, photoID)
}

// FetchPhoto the photo bytes
func FetchPhoto(db Connection, photoID int, photoName string) (*[]byte, error) {

	results, err := db.Query(fmt.Sprintf(`
		SELECT
			IMAGE
		FROM %v 
		WHERE
			ID = $1 AND
			FILENAME = $2
	`, table), photoID, photoName)

	if err != nil {
		log.Fatal("Error retrieving item")
		return nil, err
	}

	results.Next()
	var photo []byte
	if err = results.Scan(&photo); err != nil {
		return nil, err
	}

	return &photo, nil
}

// FetchThumb serve the thumbnail
func FetchThumb(db Connection, photoID int, photoName string) (*[]byte, error) {

	results, err := db.Query(fmt.Sprintf(`
		SELECT
			THUMB
		FROM %v 
		WHERE
			ID = $1 AND
			THUMBFILENAME = $2
	`, table), photoID, photoName)

	if err != nil {
		log.Fatal("Error retrieving item")
		return nil, err
	}

	results.Next()
	var thumb []byte
	if err = results.Scan(&thumb); err != nil {
		return nil, err
	}

	return &thumb, nil
}

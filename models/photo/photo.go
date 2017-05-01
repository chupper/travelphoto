package photo

import (
	"fmt"
	"log"

	"github.com/chupper/travelphoto/shared/database"
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

// Create adds a new photo for the gallery
func Create(db database.Connection, name string, description string, galleryID int) int {
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
func GetPhotos(db database.Connection, galleryID int) (*[]Photo, error) {

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
	defer results.Close()

	if err != nil {
		log.Println(err)
	}

	// not sure if this is a golang way?
	var items []Photo
	for results.Next() {
		var id, galleryID int
		var name, description, fileName, thumbfilename string
		err = results.Scan(&id, &galleryID, &name, &description, &fileName, &thumbfilename)

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
func UpdatePhoto(db database.Connection, photoID int, fileName string, photoBytes *[]byte) {
	db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			IMAGE = $1,
			FILENAME = $2
		WHERE ID = $3
		`, table), photoBytes, fileName, photoID)
}

// UpdateThumb updates the thumbnail
func UpdateThumb(db database.Connection, photoID int, fileName string, photoBytes *[]byte) {
	db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			Thumb = $1,
			ThumbFileName = $2
		WHERE ID = $3
		`, table), photoBytes, fileName, photoID)
}

// FetchPhoto the photo bytes
func FetchPhoto(db database.Connection, photoID int, photoName string) (*[]byte, error) {

	results, err := db.Query(fmt.Sprintf(`
		SELECT
			IMAGE
		FROM %v 
		WHERE
			ID = $1 AND
			FILENAME = $2
	`, table), photoID, photoName)
	defer results.Close()

	if err != nil {
		log.Fatal("Error retrieving item:", err)
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
func FetchThumb(db database.Connection, photoID int, photoName string) (*[]byte, error) {

	results, err := db.Query(fmt.Sprintf(`
		SELECT
			THUMB
		FROM %v 
		WHERE
			ID = $1 AND
			THUMBFILENAME = $2
	`, table), photoID, photoName)
	defer results.Close()

	if err != nil {
		log.Fatal("Error retrieving item", err)
		return nil, err
	}

	results.Next()
	var thumb []byte
	if err = results.Scan(&thumb); err != nil {
		return nil, err
	}

	return &thumb, nil
}

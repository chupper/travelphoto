package gallery

import (
	"fmt"

	"log"

	"github.com/chupper/travelphoto/models/photo"
	"github.com/chupper/travelphoto/shared/database"
)

const (
	table = "gallery"
)

// Gallery model
type Gallery struct {
	ID          int
	Name        string
	Description string
	Photos      []photo.Photo
}

// Create adds a gallery
func Create(db database.Connection, name string, description string) int {
	var id int
	result := db.QueryRow(fmt.Sprintf(`
		INSERT INTO %v (name, description, show)
		VALUES ($1, $2, true)
		RETURNING ID
		`, table), name, description)
	result.Scan(&id)

	return id
}

// Select returns all the galleries
func Select(db database.Connection) (*[]Gallery, error) {

	result, err := db.Query(fmt.Sprintf(`
		SELECT ID, NAME, DESCRIPTION
		FROM %v
	`, table))
	defer result.Close()

	if err != nil {
		log.Println("Error retrieving gallery: ", err)
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

// SelectAll Query for the homepage
func SelectAll(db database.Connection) *[]Gallery {

	result, err := db.Query(fmt.Sprintf(`
		SELECT 
			g.ID, 
			g.NAME,
			g.DESCRIPTION,
			p.ID,
			p.NAME,
			p.FILENAME,
			p.THUMBFILENAME
		FROM %v g
		INNER JOIN PHOTO p on (g.ID = p.GALLERYID)
		ORDER BY g.ID
	`, table))
	defer result.Close()

	if err != nil {
		log.Fatal("Error retrieving SelectAll: ", err)
	}

	// not sure if this is a golang way?
	var items []Gallery
	morePhoto := result.Next()
	for morePhoto {
		var id, photoID int
		var name, description, photoName, photoFileName, photoThumbName string
		err = result.Scan(&id, &name, &description, &photoID, &photoName, &photoFileName, &photoThumbName)

		gallery := Gallery{
			ID:          id,
			Name:        name,
			Description: description,
		}

		photos := make([]photo.Photo, 0)
		currentGalleryID := id
		for currentGalleryID == id {
			ph := photo.Photo{
				ID:            photoID,
				Name:          photoName,
				FileName:      photoFileName,
				ThumbFileName: photoThumbName,
			}
			photos = append(photos, ph)

			// append and scan next
			morePhoto = result.Next()
			if morePhoto {
				err = result.Scan(&id, &name, &description, &photoID, &photoName, &photoFileName, &photoThumbName)
			}

			// final exit clause
			if !morePhoto {
				break
			}
		}

		// assign the photos
		gallery.Photos = photos

		items = append(items, gallery)
	}

	return &items
}

// Get returns single gallery
func Get(db database.Connection, galleryID int) (*Gallery, error) {

	result, err := db.Query(fmt.Sprintf(`
		SELECT ID, NAME, DESCRIPTION
		FROM %v
		WHERE ID = $1
	`, table), galleryID)
	defer result.Close()

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
func Update(db database.Connection, gallery Gallery) error {

	_, err := db.Exec(fmt.Sprintf(`
		UPDATE %v SET
			DESCRIPTION = $1
		WHERE
			ID = $2
	`, table), gallery.Description, gallery.ID)

	return err
}

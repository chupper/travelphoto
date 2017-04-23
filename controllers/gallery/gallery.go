package gallery

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/chupper/travelphoto/controllers"
	"github.com/chupper/travelphoto/models/gallery"
	"github.com/chupper/travelphoto/models/photo"
)

// Load the gallery routes
func Load(r *mux.Router) {
	r.HandleFunc("/gallery", galleryHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/gallery/create", createGalleryHandler).Methods(http.MethodGet)
	r.HandleFunc("/gallery/{galleryid:[0-9]+}", editGalleryHandler).Methods(http.MethodGet, http.MethodPost)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {

	db := controllers.DbConnection()

	switch r.Method {
	case http.MethodGet:

		// get the galleries
		galleries, err := gallery.Select(db)
		if err != nil {
			log.Println("Error retrieving galleies: ", err)
			http.NotFound(w, r)
			return
		}

		// retrieve all the items
		t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/index.tmpl")
		t.Execute(w, galleries)

	case http.MethodPost:

		// creating the new gallery and the n photos related to the gallery
		r.ParseForm()
		galleryName := r.Form["name"][0]
		galleryDescription := r.Form["description"][0]

		galleryID := gallery.Create(db, galleryName, galleryDescription)

		// create the photos for the gallery
		for i := 1; i <= len(galleryName); i++ {
			photo.Create(db, fmt.Sprintf("Photo %v", i), "", galleryID)
		}

		// after successful create we redirect to the galleries
		http.Redirect(w, r, "/gallery", 301)
	}
}

type editGallery struct {
	Gallery gallery.Gallery
	Photos  []photo.Photo
}

func editGalleryHandler(w http.ResponseWriter, r *http.Request) {

	db := controllers.DbConnection()

	var galleryID int
	galleryID, _ = strconv.Atoi(mux.Vars(r)["galleryid"])

	switch r.Method {
	case http.MethodGet:

		// populate the edit screen
		gal, errGallery := gallery.Get(db, galleryID)
		photos, errPhotos := photo.GetPhotos(db, galleryID)
		if errGallery != nil || errPhotos != nil {
			log.Println("Error retrieving gallery: ", errGallery, errPhotos)
			http.NotFound(w, r)
			return
		}

		t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/edit.tmpl")
		t.Execute(w, editGallery{
			Gallery: *gal,
			Photos:  *photos,
		})

	case http.MethodPost:

		// updates the gallery
		r.ParseForm()
		galleryName := r.Form["name"]
		galleryDescription := r.Form["description"]

		galleryUpdate := gallery.Gallery{
			ID:          galleryID,
			Name:        galleryName[0],
			Description: galleryDescription[0],
		}

		gallery.Update(db, galleryUpdate)

		// redirect back to gallery
		http.Redirect(w, r, fmt.Sprint("/gallery/", galleryID), 301)
	}
}

func createGalleryHandler(w http.ResponseWriter, r *http.Request) {
	// return the gallery
	t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/create.tmpl")
	t.Execute(w, controllers.Page{})
}

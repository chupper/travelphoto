package gallery

import (
	"html/template"
	"log"
	"net/http"

	"database/sql"

	"github.com/chupper/travelphoto/controllers"
	"github.com/chupper/travelphoto/models/gallery"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Load the gallery routes
func Load(r *mux.Router) {
	r.HandleFunc("/gallery", galleryHandler)
	r.HandleFunc("/gallery/create", createGalleryHandler)
	r.HandleFunc("/gallery/{id}", getGalleryHandler)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Connect fail: ", err)
	}

	if r.Method == "GET" {

		// get the galleries
		galleries, err := gallery.Select(db)
		if err != nil {
			log.Fatal("Retrieve fail: ", err)
		}

		// retrieve all the items
		t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/index.tmpl")
		t.Execute(w, galleries)

	} else if r.Method == "POST" {

		// creating the new gallery
		r.ParseForm()
		galleryName := r.Form["name"][0]
		galleryDescription := r.Form["description"][0]

		if _, err := gallery.Create(db, galleryName, galleryDescription); err != nil {
			log.Fatal("Create fail: ", err)
		}

		// after successful create we redirect to the galleries
		http.Redirect(w, r, "/gallery", 301)
	}
}

func getGalleryHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/edit.tmpl")
	t.Execute(w, controllers.Page{})
}

func createGalleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// return the gallery
		t, _ := template.ParseFiles("views/base.tmpl", "views/gallery/create.tmpl")
		t.Execute(w, controllers.Page{})
	} else if r.Method == "PUT" {
		// updates the gallery
		r.ParseForm()
		galleryID := r.Form["id"]
		galleryName := r.Form["name"]
		galleryDescription := r.Form["description"]
		s := galleryName[0] + galleryDescription[0] + galleryID[0]

		w.Write([]byte(s))
	}
}

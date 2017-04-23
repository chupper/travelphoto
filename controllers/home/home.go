package home

import (
	"html/template"
	"log"
	"net/http"

	"github.com/chupper/travelphoto/controllers"
	"github.com/chupper/travelphoto/models/gallery"
	"github.com/chupper/travelphoto/models/photo"

	"github.com/gorilla/mux"
)

//Load loads the routes for the home page
func Load(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
}

type galleryView struct {
	Gallery gallery.Gallery
	Photos  []photo.Photo
}

type homePage struct {
	Galleries []galleryView
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	db, err := controllers.DbConnection()
	if err != nil {
		log.Fatal("Error Initialising Database ", err)
		return
	}

	galleries, _ := gallery.SelectAll(db)

	t, _ := template.ParseFiles("views/home/home.tmpl")
	t.Execute(w, galleries)
}

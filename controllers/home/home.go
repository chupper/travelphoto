package home

import (
	"html/template"
	"net/http"

	"github.com/chupper/travelphoto/models/gallery"
	"github.com/chupper/travelphoto/models/photo"
	"github.com/chupper/travelphoto/shared/database"

	"github.com/gorilla/mux"
)

//Load loads the routes for the home page
func Load(r *mux.Router) {
	r.HandleFunc("/", homeHandler).Methods(http.MethodGet)
}

type galleryView struct {
	Gallery gallery.Gallery
	Photos  []photo.Photo
}

type homePage struct {
	Galleries []galleryView
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	db := database.DbConnection()
	galleries := gallery.SelectAll(db)

	t, _ := template.ParseFiles("views/home/home.tmpl")
	t.Execute(w, galleries)
}

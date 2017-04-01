package home

import (
	"html/template"
	"net/http"

	"github.com/chupper/travelphoto/controllers"

	"github.com/gorilla/mux"
)

//Load loads the routes for the home page
func Load(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/home.html")
	t.Execute(w, controllers.Page{
		Title: "test",
	})
}

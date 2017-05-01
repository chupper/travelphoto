package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chupper/travelphoto/controllers/gallery"
	"github.com/chupper/travelphoto/controllers/home"
	"github.com/chupper/travelphoto/controllers/login"
	"github.com/chupper/travelphoto/controllers/photo"
	"github.com/chupper/travelphoto/shared/database"
	"github.com/chupper/travelphoto/shared/session"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	// set up the session store
	session.Configure(session.Session{
		Options: sessions.Options{
			MaxAge: 1200,
			//Secure:   true,
			HttpOnly: true,
		},
		Name:      "TravelPhoto",
		SecretKey: "thisistopsecret",
	})

	// setup the database
	database.Configure(database.DbConfig{})

	// set up routes
	r := mux.NewRouter()

	// loading the routes
	home.Load(r)
	gallery.Load(r)
	photo.Load(r)
	login.Load(r)

	// serve the static folder
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	// start the server
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println(srv.ListenAndServe())
}

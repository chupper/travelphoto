package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chupper/travelphoto/controllers/home"

	"github.com/gorilla/mux"
)

func main() {
	// start the webserver
	r := mux.NewRouter()

	// loading the routes
	home.Load(r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

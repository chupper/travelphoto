package photo

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chupper/travelphoto/controllers"
	"github.com/chupper/travelphoto/models/photo"
	"github.com/gorilla/mux"
)

// Load loads the routes
func Load(r *mux.Router) {
	r.HandleFunc("/galleryphoto/{galleryid:[0-9]+}/{photoname:[A-Z|a-z|0-9|.|_]+}", servePhoto)
	r.HandleFunc("/photo/{photoid:[0-9]+}", editPhoto)
}

func editPhoto(w http.ResponseWriter, r *http.Request) {

	log.Println("Updating Photo")

	db, err := controllers.DbConnection()
	if err != nil {
		return
	}

	var photoID int
	photoID, _ = strconv.Atoi(mux.Vars(r)["photoid"])

	// uploading new photo
	if r.Method == "POST" {

		// try for the update
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")

		if err != nil {
			log.Fatal("Not good, we couldn't read the file", handler.Filename, handler.Header)
			return
		}

		defer file.Close()
		var buff bytes.Buffer
		buff.ReadFrom(file)
		fileBytes := buff.Bytes()

		photo.UpdatePhoto(db, photoID, handler.Filename, &fileBytes)
		http.Redirect(w, r, fmt.Sprint("/gallery"), 301)
	} else {
		log.Fatal("We'll crap...")
	}
}

func servePhoto(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Photo")

	db, err := controllers.DbConnection()
	if err != nil {
		return
	}

	if r.Method == "GET" {

		galleryID, _ := strconv.Atoi(mux.Vars(r)["galleryid"])
		photoName, _ := mux.Vars(r)["photoname"]
		log.Println("Serving ", galleryID, " ", photoName)

		photoBytes, err := photo.FetchPhoto(db, galleryID, photoName)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(*photoBytes)))
		w.Write(*photoBytes)
	}
}

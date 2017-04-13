package photo

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"log"

	"bytes"

	"github.com/chupper/travelphoto/models/photo"
	"github.com/gorilla/mux"
)

// Load loads the routes
func Load(r *mux.Router) {
	r.HandleFunc("/galleryphoto/{galleryid:[0-9]+}/{photoname:[a-z0-9]}", servePhoto)
	r.HandleFunc("/photo/{photoid:[0-9]+}", editPhoto)
}

func editPhoto(w http.ResponseWriter, r *http.Request) {

	log.Println("Updating Photo")

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Connect fail: ", err)
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

		photo.UpdatePhoto(db, photoID, &fileBytes)
		http.Redirect(w, r, fmt.Sprint("/gallery"), 301)
	} else {
		log.Fatal("We'll crap...")
	}
}

func servePhoto(w http.ResponseWriter, r *http.Request) {

}

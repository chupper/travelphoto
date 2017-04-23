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
	r.HandleFunc("/galleryphoto/{photoid:[0-9]+}/{photoname:[A-Z|a-z|0-9|.|_]+}", servePhoto).Methods(http.MethodGet)
	r.HandleFunc("/gallerythumb/{photoid:[0-9]+}/{photoname:[A-Z|a-z|0-9|.|_]+}", serveThumbnail).Methods(http.MethodGet)
	r.HandleFunc("/photo/{photoid:[0-9]+}", editPhoto).Methods(http.MethodPost)
}

func editPhoto(w http.ResponseWriter, r *http.Request) {

	db := controllers.DbConnection()

	var photoID int
	photoID, _ = strconv.Atoi(mux.Vars(r)["photoid"])

	// uploading new photo
	// try for the update
	r.ParseMultipartForm(32 << 20)
	galleryID := r.FormValue("galleryid")

	// update photo if exists
	fileBytes, fileName, _ := readImage(r, "image")
	if fileBytes != nil {
		photo.UpdatePhoto(db, photoID, fileName, &fileBytes)
	}

	// update thumbnail
	fileBytes, fileName, _ = readImage(r, "thumb")
	if fileBytes != nil {
		photo.UpdateThumb(db, photoID, fileName, &fileBytes)
	}

	http.Redirect(w, r, fmt.Sprint("/gallery/", galleryID), 301)
}

func readImage(r *http.Request, name string) ([]byte, string, error) {

	file, handler, err := r.FormFile(name)

	if err != nil {
		log.Println("File not available.")
		return nil, "", nil
	}

	defer file.Close()
	var buff bytes.Buffer
	buff.ReadFrom(file)
	fileBytes := buff.Bytes()

	return fileBytes, handler.Filename, nil
}

func servePhoto(w http.ResponseWriter, r *http.Request) {

	db := controllers.DbConnection()

	photoID, _ := strconv.Atoi(mux.Vars(r)["photoid"])
	photoName, _ := mux.Vars(r)["photoname"]

	photoBytes, err := photo.FetchPhoto(db, photoID, photoName)
	if err != nil {
		log.Println("Error Fetching Photo: ", err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(*photoBytes)))
	w.Write(*photoBytes)
}

func serveThumbnail(w http.ResponseWriter, r *http.Request) {

	db := controllers.DbConnection()

	photoID, _ := strconv.Atoi(mux.Vars(r)["photoid"])
	photoName, _ := mux.Vars(r)["photoname"]

	photoBytes, err := photo.FetchThumb(db, photoID, photoName)
	if err != nil {
		log.Println("Error Fetching Photo: ", err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(*photoBytes)))
	w.Write(*photoBytes)
}

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
	r.HandleFunc("/galleryphoto/{photoid:[0-9]+}/{photoname:[A-Z|a-z|0-9|.|_]+}", servePhoto)
	r.HandleFunc("/gallerythumb/{photoid:[0-9]+}/{photoname:[A-Z|a-z|0-9|.|_]+}", serveThumbnail)
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

		// update photo if exists
		fileBytes, fileName, _ := readImage(r, "image")
		if fileBytes != nil {
			photo.UpdatePhoto(db, photoID, fileName, &fileBytes)
		}

		// update thumbnail
		thumbfileBytes, _, _ := readImage(r, "thumb")
		if thumbfileBytes != nil {
			photo.UpdateThumb(db, photoID, &thumbfileBytes)
		}

		http.Redirect(w, r, fmt.Sprint("/gallery"), 301)
	} else {
		log.Fatal("We'll crap...")
	}
}

func readImage(r *http.Request, name string) ([]byte, string, error) {

	val, ok := r.Form[name]
	if ok || len(val) != 0 {
		return nil, "", nil
	}

	file, handler, err := r.FormFile(name)

	if err != nil {
		log.Fatal("Error reading file: ", name, " ", err)
		return nil, "", nil
	}

	defer file.Close()
	var buff bytes.Buffer
	buff.ReadFrom(file)
	fileBytes := buff.Bytes()

	return fileBytes, handler.Filename, nil
}

func servePhoto(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Photo")

	db, err := controllers.DbConnection()
	if err != nil {
		return
	}

	if r.Method == "GET" {

		photoID, _ := strconv.Atoi(mux.Vars(r)["photoid"])
		photoName, _ := mux.Vars(r)["photoname"]
		log.Println("Serving ", photoID, " ", photoName)

		photoBytes, err := photo.FetchPhoto(db, photoID, photoName)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(*photoBytes)))
		w.Write(*photoBytes)
	}
}

func serveThumbnail(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Thumb")

	db, err := controllers.DbConnection()
	if err != nil {
		return
	}

	if r.Method == "GET" {

		photoID, _ := strconv.Atoi(mux.Vars(r)["photoid"])
		photoName, _ := mux.Vars(r)["photoname"]
		log.Println("Serving ", photoID, " ", photoName)

		photoBytes, err := photo.FetchThumb(db, photoID, photoName)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(*photoBytes)))
		w.Write(*photoBytes)
	}
}

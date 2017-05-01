package login

import (
	"html/template"
	"net/http"

	"log"

	"github.com/chupper/travelphoto/models/user"
	"github.com/chupper/travelphoto/shared/database"
	"github.com/chupper/travelphoto/shared/session"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Load loads the hanlders for login / logout
func Load(r *mux.Router) {
	r.HandleFunc("/login", loginPage).Methods(http.MethodGet)
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler).Methods(http.MethodGet)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/base.tmpl", "views/login.tmpl")
	t.Execute(w, nil)
}

// we'll map this with gorilla later
type loginForm struct {
	UserName string
	Password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	db := database.DbConnection()

	// get the user name and the password
	r.ParseForm()
	userName := r.PostFormValue("username")
	password := r.PostFormValue("password")

	isValid, err := user.Validate(db, userName, password)
	if err != nil {
		log.Println("Error logging in:", err.Error())
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}

	session := session.Instance(r)
	if isValid {
		log.Println(userName, password, session.Values["is_authenticated"])
		session.Values["is_authenticated"] = true
		sessions.Save(r, w)

		http.Redirect(w, r, "/gallery", http.StatusSeeOther)
	} else {
		// build a flash message here
		session.AddFlash("Validation", "Error: Username Password is invalid")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	// do whatever we need to do to logout
	log.Println("Logging out user")
	s := session.Instance(r)
	s.Options.MaxAge = -1
	sessions.Save(r, w)

	// redirect back to main page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

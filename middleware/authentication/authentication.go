package authentication

import (
	"log"
	"net/http"

	"github.com/chupper/travelphoto/shared/session"
)

// Authenticated middleware
// redirects to login page if try to access a non logged in page when
// not authenticated
func Authenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we should load the user into the session if this is active?
		s := session.Instance(r)

		log.Println("Checking Authenticated")
		if s != nil && s.Values["is_authenticated"] == true {
			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
	})
}

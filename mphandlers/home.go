package mphandlers

import (
	"net/http"
)

// handleHome handles HTTP requests for the homepage.
func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html", nil)
}

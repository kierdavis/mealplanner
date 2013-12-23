package mphandlers

import (
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html", nil)
}

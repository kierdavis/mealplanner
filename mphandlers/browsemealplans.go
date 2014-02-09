package mphandlers

import (
	"net/http"
)

func handleBrowseMealPlans(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "browse-mps.html", nil)
}

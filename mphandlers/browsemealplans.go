package mphandlers

import (
	"net/http"
)

func handleBrowseMealPlans(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "browse-meal-plans.html", nil)
}

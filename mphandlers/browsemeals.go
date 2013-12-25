package mphandlers

import (
	//"database/sql"
	//"github.com/kierdavis/mealplanner/mpdata"
	//"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

// handleBrowseMeals handles HTTP requests for the meal list.
func handleBrowseMeals(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "browse-meals.html", nil)
}

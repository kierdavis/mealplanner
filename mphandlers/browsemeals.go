package mphandlers

import (
	//"database/sql"
	//"github.com/kierdavis/mealplanner/mpdata"
	//"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

func handleBrowseMeals(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "browse-meals.html", nil)
}

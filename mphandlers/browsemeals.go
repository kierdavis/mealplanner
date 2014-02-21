package mphandlers

import (
	//"database/sql"
	//"github.com/kierdavis/mealplanner/mpdata"
	//"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
	"strconv"
)

// highlightData contains information passed to the meal browser template
// regarding which meal, if any, should be highlighted in the list.
type highlightData struct {
	Highlight bool
	MealID    uint64
}

// handleBrowseMeals handles HTTP requests for the meal list.
func handleBrowseMeals(w http.ResponseWriter, r *http.Request) {
	var hd highlightData

	highlightStr := r.FormValue("highlight")
	if highlightStr != "" {
		mealID, err := strconv.ParseUint(highlightStr, 10, 64)
		if err != nil {
			httpError(w, BadRequestError)
			return
		}

		hd.Highlight = true
		hd.MealID = mealID
	}

	renderTemplate(w, "browse-meals.html", hd)
}

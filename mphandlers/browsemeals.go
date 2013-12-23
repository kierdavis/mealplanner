package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

func handleBrowseMeals(w http.ResponseWriter, r *http.Request) {
	var mts []mpdata.MealWithTags
	
	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		mts, err = mpdb.ListMealsWithTags(db, true)
		return err
	})
	if err != nil {
		serverError(w, err)
		return
	}
	
	renderTemplate(w, "browse-meals.html", mts)
}

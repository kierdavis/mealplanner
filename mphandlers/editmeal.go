package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

// handleEditMealForm handles HTTP requests for the "edit meal" form.
func handleEditMealForm(w http.ResponseWriter, r *http.Request) {
	mealID, ok := getUint64Var(r, "mealid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}

	var mt mpdata.MealWithTags

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		mt, err = mpdb.GetMealWithTags(db, mealID)
		return err
	})

	if err != nil {
		serverError(w, err)
		return
	}

	renderTemplate(w, "edit-meal-form.html", mt)
}

package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

func handleEditMealForm(w http.ResponseWriter, r *http.Request) {
	mealID, ok := getUint64Var(r, "mealid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}
	
	var meal *mpdata.Meal
	
	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		meal, err = mpdb.GetMealWithTags(db, mealID)
		return err
	})
	
	if err != nil {
		serverError(w, err)
		return
	}
	
	renderTemplate(w, "edit-meal-form.html", meal)
}

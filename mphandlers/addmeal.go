package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

// handleAddMealForm handles HTTP requests for the "new meal" form.
func handleAddMealForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "edit-meal-form.html", nil)
}

// handleAddMealAction handles HTTP requests for submission of the "new meal"
// form.
func handleAddMealAction(w http.ResponseWriter, r *http.Request) {
	// Parse the POST request body
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	// Create a MealWithTags value from the form fields
	mt := mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			Name:      r.FormValue("name"),
			RecipeURL: r.FormValue("recipe"),
			Favourite: r.FormValue("favourite") != "",
		},
		Tags: r.Form["tags"],
	}

	// Add the record to the database
	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.AddMealWithTags(tx, mt)
		})
	})

	if err != nil {
		serverError(w, err)
		return
	}

	// Redirect to list of meals
	redirect(w, http.StatusSeeOther, "/meals")
}

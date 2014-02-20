package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
	"strconv"
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
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mt, err = mpdb.GetMealWithTags(tx, mealID)
			return err
		})
	})

	if err != nil {
		serverError(w, err)
		return
	}

	if mt.Meal == nil {
		httpError(w, NotFoundError)
		return
	}

	renderTemplate(w, "edit-meal-form.html", mt)
}

// handleEditMealAction handles HTTP requests for submission of the "edit meal"
// form.
func handleEditMealAction(w http.ResponseWriter, r *http.Request) {
	// Get the meal ID from the URL
	mealID, ok := getUint64Var(r, "mealid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}

	// Parse the POST request body
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	// Create a MealWithTags value from the form fields
	mt := mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			ID:        mealID,
			Name:      r.FormValue("name"),
			RecipeURL: r.FormValue("recipe"),
			Favourite: r.FormValue("favourite") != "",
		},
		Tags: r.Form["tags"],
	}

	// Update the database record
	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.UpdateMealWithTags(tx, mt)
		})
	})
	if err != nil {
		serverError(w, err)
		return
	}

	// Redirect to list of meals
	redirect(w, http.StatusSeeOther, "/meals?highlight="+strconv.FormatUint(mealID, 10))
}

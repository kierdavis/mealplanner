package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
	"strconv"
	"time"
)

// handleCreateMealPlanForm handles HTTP requests for the meal plan creation
// form.
func handleCreateMealPlanForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "create-mp-form.html", nil)
}

// handleCreateMealPlanAction handles HTTP requests for the submission of the
// meal plan creation form.
func handleCreateMealPlanAction(w http.ResponseWriter, r *http.Request) {
	// Parse the POST request body
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	startDate, err := time.Parse(mpdata.DatepickerDateFormat, r.FormValue("start"))
	if err != nil {
		httpError(w, BadRequestError)
		return
	}

	endDate, err := time.Parse(mpdata.DatepickerDateFormat, r.FormValue("end"))
	if err != nil {
		httpError(w, BadRequestError)
		return
	}

	if startDate.After(endDate) {
		httpError(w, BadRequestError)
		return
	}

	auto := r.FormValue("auto") == "true"

	// Create a MealPlan object
	mp := &mpdata.MealPlan{
		StartDate: startDate,
		EndDate:   endDate,
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			// Add mp to the database
			err = mpdb.AddMealPlan(tx, mp)
			if err != nil {
				return err
			}

			// Optionally fill the meal plan automatically
			if auto {
				err = mpdb.AutoFillMealPlan(tx, mp)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
	if err != nil {
		serverError(w, err)
		return
	}

	redirect(w, http.StatusSeeOther, "/mealplans/"+strconv.FormatUint(mp.ID, 10)+"/edit")
}

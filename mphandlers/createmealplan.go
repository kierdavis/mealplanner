package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
	"strconv"
	"time"
)

func handleCreateMealPlanForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "create-mp-form.html", nil)
}

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

	// Create a MealPlan
	mp := &mpdata.MealPlan{
		StartDate: startDate,
		EndDate:   endDate,
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.AddMealPlan(tx, mp)
		})
	})
	if err != nil {
		serverError(w, err)
		return
	}

	redirect(w, http.StatusSeeOther, "/mealplans/"+strconv.FormatUint(mp.ID, 10))
}

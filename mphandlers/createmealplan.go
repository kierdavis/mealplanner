package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
	"strconv"
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
	
	// Create a MealPlan
	mp := &mpdata.MealPlan{
		Start: time.Parse(mpdata.DatepickerFormat, r.FormValue("start")),
		End: time.Parse(mpdata.DatepickerFormat, r.FormValue("end")),
	}
	
	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.CreateMealPlan(tx, mp)
		})
	})
	if err != nil {
		serverError(w, err)
		return
	}
	
	redirect(w, http.StatusSeeOther, "/mealplans/" + strconv.FormatUint(mp.ID, 10))
}

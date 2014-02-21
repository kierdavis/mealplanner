package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

// deleteMPData contains information passed to the meal plan deletion template
// regarding the meal plan that is being deleted.
type deleteMPData struct {
	MP          *mpdata.MealPlan
	NumServings int
}

// handleDeleteMealPlanForm handles HTTP requests for the meal plan deletion
// confirmation page.
func handleDeleteMealPlanForm(w http.ResponseWriter, r *http.Request) {
	mpID, ok := getUint64Var(r, "mealplanid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}

	var mp *mpdata.MealPlan
	var numServings int

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mp, err = mpdb.GetMealPlan(tx, mpID)
			if err != nil {
				return err
			}

			numServings, err = mpdb.CountServings(tx, mpID)
			return err
		})
	})

	if err != nil {
		serverError(w, err)
		return
	}

	if mp == nil {
		httpError(w, NotFoundError)
		return
	}

	renderTemplate(w, "delete-mp-form.html", deleteMPData{mp, numServings})
}

// handleDeleteMealPlanAction handles HTTP requests for submission of the
// meal plan deletion form.
func handleDeleteMealPlanAction(w http.ResponseWriter, r *http.Request) {
	mpID, ok := getUint64Var(r, "mealplanid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			err = mpdb.DeleteServings(tx, mpID)
			if err != nil {
				return err
			}

			return mpdb.DeleteMealPlan(tx, mpID)
		})
	})

	if err != nil {
		serverError(w, err)
		return
	}

	redirect(w, http.StatusSeeOther, "/mealplans")
}

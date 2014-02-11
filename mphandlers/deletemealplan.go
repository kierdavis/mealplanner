package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

type deleteMPFormValues struct {
	MP *mpdata.MealPlan
	NumServings int
}

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
	
	renderTemplate(w, "delete-mp-form.html", deleteMPFormValues{mp, numServings})
}

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

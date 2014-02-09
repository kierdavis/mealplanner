package mphandlers

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/http"
)

func handleViewMealPlan(w http.ResponseWriter, r *http.Request) {
	mpID, ok := getUint64Var(r, "mealplanid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}

	var mp *mpdata.MealPlan

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mp, err = mpdb.GetMealPlan(tx, mpID)
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

	renderTemplate(w, "view-mp.html", mp)
}

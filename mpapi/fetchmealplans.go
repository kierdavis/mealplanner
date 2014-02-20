package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"time"
)

// fetchMealPlans handles an API call to return a list of meal plans that
// overlap with a specified inclusive date range. Expected parameters: from, to.
// Returns: an array of meal plan objects.
func fetchMealPlans(params url.Values) (response JSONResponse) {
	from, err := time.Parse(mpdata.JsonDateFormat, params.Get("from"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'from' parameter"}
	}

	to, err := time.Parse(mpdata.JsonDateFormat, params.Get("to"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'to' parameter"}
	}

	var mps []*mpdata.MealPlan

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mps, err = mpdb.ListMealPlansBetween(tx, from, to)
			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: mps}
}

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

func fetchMealPlans(params url.Values) (response JsonResponse) {
	from, err := time.Parse(mpdata.JsonDateFormat, params.Get("from"))
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'from' parameter"}
	}
	
	to, err := time.Parse(mpdata.JsonDateFormat, params.Get("to"))
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'to' parameter"}
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
		return JsonResponse{Error: "Database error"}
	}
	
	return JsonResponse{Success: mps}
}

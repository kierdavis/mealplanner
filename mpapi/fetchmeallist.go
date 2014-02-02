package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
)

// fetchMealList handles an API call to fetch a list of all meals in the
// database.
func fetchMealList(params url.Values) (response JsonResponse) {
	var mts []mpdata.MealWithTags

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mts, err = mpdb.ListMealsWithTags(tx, true)
			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: mts}
}

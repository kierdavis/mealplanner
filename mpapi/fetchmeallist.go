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
// database. Expected parameters: none. Returns: an array of meal/tags objects.
func fetchMealList(params url.Values) (response JSONResponse) {
	var mts []mpdata.MealWithTags

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mts, err = mpdb.ListMealsWithTags(tx, true)
			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: mts}
}

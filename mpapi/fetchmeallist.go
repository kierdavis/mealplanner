package mpapi

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"log"
	"net/url"
	"regexp"
)

var wordRegexp = regexp.MustCompile("\\w+")

// fetchMealList handles an API call to fetch a list of all meals in the
// database. Expected parameters: none. Returns: an array of meal/tags objects.
func fetchMealList(params url.Values) (response JSONResponse) {
	query := params.Get("query")
	var words []string

	if query != "" {
		words = wordRegexp.FindAllString(query, -1)
	}

	var mts []mpdata.MealWithTags

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			if query == "" {
				mts, err = mpdb.ListMealsWithTags(tx, true)
			} else {
				mts, err = mpdb.SearchMealsWithTags(tx, words, true)
			}

			return err
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: mts}
}

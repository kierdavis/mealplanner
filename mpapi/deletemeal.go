package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"strconv"
)

// deleteMeal handles an API call to delete a meal.
func deleteMeal(params url.Values) (response JsonResponse) {
	mealID, err := strconv.ParseUint(params.Get("mealid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealid' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			err = mpdb.DeleteServingsOf(tx, mealID)
			if err != nil {
				return err
			}
			
			return mpdb.DeleteMealWithTags(tx, mealID)
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: nil}
}

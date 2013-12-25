package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
)

func fetchMealList(params url.Values) (response JsonResponse) {
	var mts []mpdata.MealWithTags
	
	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		mts, err = mpdb.ListMealsWithTags(db, true)
		return err
	})
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}
	
	return JsonResponse{Success: mts}
}
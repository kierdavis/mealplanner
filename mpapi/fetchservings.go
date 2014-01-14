package mpapi

import (
	//"github.com/kierdavis/mealplanner/mpdata"
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"strconv"
	"time"
)

type fetchServingsRecord struct {
	Date     time.Time `json:"date"`
	HasMeal  bool      `json:"hasmeal"`
	MealID   uint64    `json:"mealid"`
	MealName string    `json:"mealname"`
}

func fetchServings(params url.Values) (response JsonResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	var results []*fetchServingsRecord

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			mps, err := mpdb.GetMealPlanWithServings(tx, mpID)
			if err != nil {
				return err
			}

			for _, date := range mps.MealPlan.Days() {
				ts := &fetchServingsRecord{
					Date: date,
				}

				for _, serving := range mps.Servings {
					if serving.Date == date {
						ts.HasMeal = true
						ts.MealID = serving.MealID

						meal, err := mpdb.GetMeal(tx, serving.MealID)
						if err != nil {
							return err
						}
						
						if meal == nil {
							fmt.Fprintf(os.Stderr, "Warning: meal plan %d -> serving %s points to nonexistent meal %d\n", mpID, date.Format("2006-01-02"), serving.MealID)
							ts.MealName = "???"
						} else {
							ts.MealName = meal.Name
						}

						break
					}
				}

				results = append(results, ts)
			}

			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: results}
}

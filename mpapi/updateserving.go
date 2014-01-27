package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"strconv"
	"time"
)

func updateServing(params url.Values) (response JsonResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	dateServed, err := time.Parse(mpdata.JsonDateFormat, params.Get("date"))
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'date' parameter"}
	}

	mealID, err := strconv.ParseUint(params.Get("mealid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealid' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			err = mpdb.DeleteServing(tx, mpID, dateServed)
			if err != nil {
				return err
			}

			s := &mpdata.Serving{
				MealPlanID: mpID,
				Date:       dateServed,
				MealID:     mealID,
			}

			return mpdb.AddServing(tx, s)
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: nil}
}

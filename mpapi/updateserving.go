package mpapi

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"log"
	"net/url"
	"strconv"
	"time"
)

// updateServing implements an API call to update a meal serving for a meal
// plan with a new meal ID, removing the old serving if it already exists.
// Expected parameters: mealplanid, date, mealid. Returns: nothing.
func updateServing(params url.Values) (response JSONResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	dateServed, err := time.Parse(mpdata.JSONDateFormat, params.Get("date"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'date' parameter"}
	}

	mealID, err := strconv.ParseUint(params.Get("mealid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealid' parameter"}
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
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: nil}
}

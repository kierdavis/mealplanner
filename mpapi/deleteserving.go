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

// deleteServing handles an API call to delete a meal serving. Expected
// parameters: mealplanid, date. Returns: nothing.
func deleteServing(params url.Values) (response JSONResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	servingDate, err := time.Parse(mpdata.JSONDateFormat, params.Get("date"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'date' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.DeleteServing(tx, mpID, servingDate)
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: nil}
}

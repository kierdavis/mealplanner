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

// fetchSuggestions handles an API call to generate suggestions for a given date.
// Expected parameters: date. Returns: an array of suggestion objects.
func fetchSuggestions(params url.Values) (response JSONResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}
	
	dateServed, err := time.Parse(mpdata.JSONDateFormat, params.Get("date"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'date' parameter"}
	}

	var suggs []*mpdata.Suggestion

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			suggs, err = mpdb.GenerateSuggestions(tx, mpID, dateServed)
			return err
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: suggs}
}

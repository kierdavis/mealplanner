package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"time"
)

// fetchSuggestions handles an API call to generate suggestions for a given date.
// Expected parameters: date. Returns: an array of suggestion objects.
func fetchSuggestions(params url.Values) (response JSONResponse) {
	dateServed, err := time.Parse(mpdata.JSONDateFormat, params.Get("date"))
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'date' parameter"}
	}

	var suggs []*mpdata.Suggestion

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			suggs, err = mpdb.GenerateSuggestions(tx, dateServed)
			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: suggs}
}

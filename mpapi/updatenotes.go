package mpapi

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdb"
	"log"
	"net/url"
	"strconv"
)

// updateNotes implements an API call to update the notes associated with a
// meal plan. Expected parameters: mealplanid, notes. Returns: nothing.
func updateNotes(params url.Values) (response JSONResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	notes := params.Get("notes")
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'notes' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.UpdateNotes(tx, mpID, notes)
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: nil}
}

package mpapi

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdb"
	"log"
	"net/url"
	"strconv"
)

// toggleFavourite implements an API call to toggle the "favourite" status of
// a given meal. Expected paramaters: mealid. Returns: the updated "favourite"
// status of the meal.
func toggleFavourite(params url.Values) (response JSONResponse) {
	mealID, err := strconv.ParseUint(params.Get("mealid"), 10, 64)
	if err != nil {
		return JSONResponse{Error: "Invalid or missing 'mealid' parameter"}
	}

	var isFavourite bool

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			isFavourite, err = mpdb.ToggleFavourite(tx, mealID)
			return err
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: isFavourite}
}

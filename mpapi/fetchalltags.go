package mpapi

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdb"
	"log"
	"net/url"
)

// fetchAllTags handles an API call to obtain a list of all tags present in the
// database, without duplicates and in alphabetical order. Expected parameters:
// none. Returns: an array of tags.
func fetchAllTags(params url.Values) (response JSONResponse) {
	var tags []string

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			tags, err = mpdb.ListAllTags(tx, true)
			return err
		})
	})

	if err != nil {
		log.Printf("Database error: %s\n", err.Error())
		return JSONResponse{Error: "Database error"}
	}

	return JSONResponse{Success: tags}
}

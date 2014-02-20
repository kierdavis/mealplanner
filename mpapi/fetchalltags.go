package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
)

// fetchAllTags handles an API call to obtain a list of all tags present in the
// database, without duplicates and in alphabetical order. Expected parameters:
// none. Returns: an array of tags.
func fetchAllTags(params url.Values) (response JsonResponse) {
	var tags []string

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			tags, err = mpdb.ListAllTags(tx, true)
			return err
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: tags}
}

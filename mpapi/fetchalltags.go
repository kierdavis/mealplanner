package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
)

// fetchAllTags handles an API call to obtain a list of all tags present in the
// database, without duplicates and in alphabetical order.
func fetchAllTags(params url.Values) (response JsonResponse) {
	var tags []string

	err := mpdb.WithConnection(func(db *sql.DB) (err error) {
		tags, err = mpdb.ListAllTags(db, true)
		return err
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: tags}
}
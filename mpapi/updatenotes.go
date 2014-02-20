package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"strconv"
)

// updateNotes implements an API call to update the notes associated with a
// meal plan. Expected parameters: mealplanid, notes. Returns: nothing.
func updateNotes(params url.Values) (response JsonResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	notes := params.Get("notes")
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'notes' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.UpdateNotes(tx, mpID, notes)
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: nil}
}

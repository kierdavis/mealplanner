package mpapi

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"github.com/kierdavis/mealplanner/mpdb"
	"net/url"
	"os"
	"strconv"
	"time"
)

// deleteServing handles an API call to delete a meal serving.
func deleteServing(params url.Values) (response JsonResponse) {
	mpID, err := strconv.ParseUint(params.Get("mealplanid"), 10, 64)
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'mealplanid' parameter"}
	}

	servingDate, err := time.Parse(mpdata.JsonDateFormat, params.Get("date"))
	if err != nil {
		return JsonResponse{Error: "Invalid or missing 'date' parameter"}
	}

	err = mpdb.WithConnection(func(db *sql.DB) (err error) {
		return mpdb.WithTransaction(db, func(tx *sql.Tx) (err error) {
			return mpdb.DeleteServing(tx, mpID, servingDate)
		})
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s\n", err.Error())
		return JsonResponse{Error: "Database error"}
	}

	return JsonResponse{Success: nil}
}

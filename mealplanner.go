// Command mealplanner is the main entry point of the application. It simply
// runs the *mux.Router provided by mphandlers.CreateMux() as an HTTP server.
package main

import (
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"github.com/kierdavis/mealplanner/mphandlers"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := mpdb.InitDB(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database Error: %s\n", err)
		os.Exit(1)
	}

	m := mphandlers.CreateMux()
	err = http.ListenAndServe("localhost:8080", m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server Error: %s\n", err)
		os.Exit(1)
	}
}

// Command mealplanner is the main entry point of the application. It simply
// runs the *mux.Router provided by mphandlers.CreateMux() as an HTTP server.
package main

import (
	"flag"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"github.com/kierdavis/mealplanner/mphandlers"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbSource = flag.String("dbsource", "", "database source, in the form USER:PASS@unix(/PATH/TO/SOCKET)/DB or USER:PASS@tcp(HOST:PORT)/DB")
	host     = flag.String("host", "", "hostname to listen on")
	port     = flag.Int("port", 8080, "port to listen on")
	debug    = flag.Bool("debug", false, "debug mode")
	testdata = flag.Bool("testdata", false, "clear the database and insert test data")
)

func main() {
	flag.Parse()

	source := *dbSource
	if source == "" {
		source = os.Getenv("MPDBSOURCE")
		if source == "" {
			fmt.Fprintf(os.Stderr, "Please specify a non-empty -dbsource flag or set the MPDBSOURCE environment variable.\n")
			os.Exit(1)
		}
	}

	mpdb.DBSource = source

	err := mpdb.InitDB(*debug, *testdata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database Error: %s\n", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%d", *host, *port)
	m := mphandlers.CreateMux()

	app := http.Handler(m)
	if *debug {
		app = mphandlers.LoggingHandler{Handler: app}

		fmt.Printf("Listening on %s\n", listenAddr)
	}

	err = http.ListenAndServe(listenAddr, app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server Error: %s\n", err)
		os.Exit(1)
	}
}

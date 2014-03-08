// Command mealplanner is the main entry point of the application. It simply
// runs the *mux.Router provided by mphandlers.CreateMux() as an HTTP server.
package main

import (
	"flag"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdb"
	"github.com/kierdavis/mealplanner/mphandlers"
	"github.com/kierdavis/mealplanner/mpresources"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbSource    = flag.String("dbsource", "", "database source, in the form USER:PASS@unix(/PATH/TO/SOCKET)/DB or USER:PASS@tcp(HOST:PORT)/DB")
	host        = flag.String("host", "", "hostname to listen on")
	port        = flag.Int("port", 8080, "port to listen on")
	debug       = flag.Bool("debug", false, "debug mode")
	testdata    = flag.Bool("testdata", false, "clear the database and insert test data")
	resourceDir = flag.String("resourcedir", "", "path to directory containing the resources used by the application")
)

func main() {
	flag.Parse()

	source := *dbSource
	if source == "" {
		source = os.Getenv("MPDBSOURCE")
		if source == "" {
			fmt.Println("Please specify a non-empty -dbsource flag or set the MPDBSOURCE environment variable.")
			os.Exit(1)
		}
	}
	
	resDir := *resourceDir
	if resDir == "" {
		resDir = os.Getenv("MPRESDIR")
	}

	mpdb.DBSource = source
	mpresources.SetResourceDir(resDir)
	mpresources.GetTemplates() // Check that the templates load correctly

	err := mpdb.InitDB(*debug, *testdata)
	if err != nil {
		log.Printf("Database error during startup: %s\n", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%d", *host, *port)
	m := mphandlers.CreateMux()

	app := http.Handler(m)
	if *debug {
		app = mphandlers.LoggingHandler{Handler: app}

		log.Printf("Listening on %s\n", listenAddr)
	}

	err = http.ListenAndServe(listenAddr, app)
	if err != nil {
		log.Printf("Server error in HTTP listener: %s\n", err)
		os.Exit(1)
	}
}

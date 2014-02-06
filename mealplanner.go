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
	host = flag.String("host", "", "hostname to listen on")
	port = flag.Int("port", 8080, "port to listen on")
	debug = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()
	
	mpdb.Source = *dbSource
	if mpdb.Source == "" {
		fmt.Fprintf(os.Stderr, "Please specify a non-empty -dbsource option.\n")
		os.Exit(1)
	}
	
	err := mpdb.InitDB(*debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database Error: %s\n", err)
		os.Exit(1)
	}

	m := mphandlers.CreateMux()
	listenAddr := fmt.Sprintf("%s:%d", *host, *port)
	
	fmt.Printf("Listening on %s\n", listenAddr)

	err = http.ListenAndServe(listenAddr, m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server Error: %s\n", err)
		os.Exit(1)
	}
}

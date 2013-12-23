package main

import (
	"fmt"
	"github.com/kierdavis/mealplanner/mphandlers"
	"net/http"
	"os"
)

func main() {
	m := mphandlers.CreateMux()
	
	err := http.ListenAndServe(":8080", m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

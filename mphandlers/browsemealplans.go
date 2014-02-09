package mphandlers

import (
	"net/http"
	"time"
)

type showingDate struct {
	Year  int
	Month int
}

func handleBrowseMealPlans(w http.ResponseWriter, r *http.Request) {
	showing := time.Now()
	showingStr := r.FormValue("showing")

	if showingStr != "" {
		var err error
		showing, err = time.Parse("2006-01-02", showingStr)
		if err != nil {
			showing = time.Now()
		}
	}

	sd := showingDate{showing.Year(), int(showing.Month() - 1)}
	renderTemplate(w, "browse-mps.html", sd)
}

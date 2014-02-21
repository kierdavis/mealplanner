package mphandlers

import (
	"net/http"
	"time"
)

// showingData contains the information passed to the meal plan browser
// template regarding which year/month to display.
type showingData struct {
	Year  int
	Month int
}

// handleBrowseMealPlans handles HTTP requests for the meal plan browser.
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

	sd := showingData{showing.Year(), int(showing.Month() - 1)}
	renderTemplate(w, "browse-mps.html", sd)
}

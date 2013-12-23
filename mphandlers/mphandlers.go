// Package mphandlers defines the HTTP handlers for the application.
package mphandlers

import (
	"github.com/gorilla/mux"
	"github.com/kierdavis/mealplanner/mpresources"
	"net/http"
)

// CreateMux creates a *mux.Router and attaches the application's HTTP handlers
// to it.
func CreateMux() (m *mux.Router) {
	m = mux.NewRouter()
	
	// Handle static files
	staticHandler := http.FileServer(http.Dir(mpresources.StaticDir))
	m.Path("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	m.Path("/").Methods("GET", "HEAD").HandlerFunc(handleHome)

	addMeal := m.Path("/meals/new").Subrouter()
	addMeal.Methods("GET", "HEAD").HandlerFunc(handleAddMealForm)
	addMeal.Methods("POST").HandlerFunc(handleAddMealAction)

	editMeal := m.Path("/meals/{mealid:[0-9]+}/edit").Subrouter()
	editMeal.Methods("GET", "HEAD").HandlerFunc(handleEditMealForm)
	//editMeal.Methods("POST").HandlerFunc(handleEditMealAction)

	return m
}

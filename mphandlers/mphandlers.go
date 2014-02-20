// Package mphandlers defines the HTTP handlers for the application.
package mphandlers

import (
	"github.com/gorilla/mux"
	"github.com/kierdavis/mealplanner/mpapi"
	"github.com/kierdavis/mealplanner/mpresources"
	"net/http"
)

// CreateMux creates a *mux.Router and attaches the application's HTTP handlers
// to it.
func CreateMux() (m *mux.Router) {
	m = mux.NewRouter()

	// Static files
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(mpresources.StaticDir)))
	m.PathPrefix("/static/").Handler(staticHandler)

	// Dynamic handlers
	m.Path("/").Methods("GET", "HEAD").HandlerFunc(handleHome)
	m.Path("/meals").Methods("GET", "HEAD").HandlerFunc(handleBrowseMeals)
	m.Path("/mealplans").Methods("GET", "HEAD").HandlerFunc(handleBrowseMealPlans)
	m.Path("/mealplans/{mealplanid:[0-9]+}").Methods("GET", "HEAD").HandlerFunc(handleViewMealPlan)
	m.Path("/mealplans/{mealplanid:[0-9]+}/edit").Methods("GET", "HEAD").HandlerFunc(handleEditMealPlan)
	m.Path("/api").Methods("POST").HandlerFunc(mpapi.HandleAPICall)

	addMeal := m.Path("/meals/new").Subrouter()
	addMeal.Methods("GET", "HEAD").HandlerFunc(handleAddMealForm)
	addMeal.Methods("POST").HandlerFunc(handleAddMealAction)

	editMeal := m.Path("/meals/{mealid:[0-9]+}/edit").Subrouter()
	editMeal.Methods("GET", "HEAD").HandlerFunc(handleEditMealForm)
	editMeal.Methods("POST").HandlerFunc(handleEditMealAction)

	createMP := m.Path("/mealplans/new").Subrouter()
	createMP.Methods("GET", "HEAD").HandlerFunc(handleCreateMealPlanForm)
	createMP.Methods("POST").HandlerFunc(handleCreateMealPlanAction)

	deleteMP := m.Path("/mealplans/{mealplanid:[0-9]+}/delete").Subrouter()
	deleteMP.Methods("GET", "HEAD").HandlerFunc(handleDeleteMealPlanForm)
	deleteMP.Methods("POST").HandlerFunc(handleDeleteMealPlanAction)

	return m
}

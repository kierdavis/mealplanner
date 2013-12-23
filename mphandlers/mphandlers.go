package mphandlers

import (
	"github.com/gorilla/mux"
)

func CreateMux() (m *mux.Router) {
	m = mux.NewRouter()
	
	m.Path("/").Methods("GET", "HEAD").HandlerFunc(handleHome)
	
	addMeal := m.Path("/meals/new").Subrouter()
	addMeal.Methods("GET", "HEAD").HandlerFunc(handleAddMealForm)
	addMeal.Methods("POST").HandlerFunc(handleAddMealAction)
	
	editMeal := m.Path("/meals/{mealid:[0-9]+}/edit").Subrouter()
	editMeal.Methods("GET", "HEAD").HandlerFunc(handleEditMealForm)
	//editMeal.Methods("POST").HandlerFunc(handleEditMealAction)
	
	return m
}

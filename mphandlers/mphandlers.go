package mphandlers

import (
	"github.com/gorilla/mux"
)

func CreateMux() (m *mux.Router) {
	m := mux.NewRouter()
	
	m.Path("").Method("GET").HandleFunc(handlehome)
	
	addMeal := m.Path("/meals/new").Subrouter()
	addMeal.Method("GET").HandleFunc(handleGetAddMeal)
	addMeal.Method("POST").HandleFunc(handlePostAddMeal)
	
	return m
}

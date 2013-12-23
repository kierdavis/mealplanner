package mphandlers

import (
	"fmt"
	"net/http"
)

func handleAddMeal(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleAddMealGET(w, r)
	case "POST":
		handleAddMealPOST(w, r)
	default:
		invalidMethod(w, "GET, POST")
	}
}

func handleAddMealGET(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "edit-meal-form.html", nil)
}

func handleAddMealPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%v\n", r.Form)
}

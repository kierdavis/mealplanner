package mphandlers

import (
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"net/http"
)

// handleAddMealForm handles HTTP requests for the "new meal" form.
func handleAddMealForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "edit-meal-form.html", nil)
}

// handleAddMealAction handles HTTP requests for submission of the "new meal"
// form.
func handleAddMealAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	mt := &mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			Name:      r.FormValue("name"),
			RecipeURL: r.FormValue("recipe"),
			Favourite: r.FormValue("favourite") != "",
		},
		Tags: r.Form["tags"],
	}

	fmt.Printf("%#v\n", mt)
}

package mphandlers

import (
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"net/http"
)

func handleAddMealForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "edit-meal-form.html", nil)
}

func handleAddMealAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}
	
	meal := &mpdata.Meal{
		Name: r.FormValue("name"),
		RecipeURL: r.FormValue("recipe"),
		Favourite: r.FormValue("favourite") != "",
		HasTags: true,
	}
	
	tags, ok := r.Form["tags"]
	if ok {
		meal.Tags = tags
	}
	
	fmt.Printf("%#v\n", meal)
}

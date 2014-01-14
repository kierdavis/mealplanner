package mphandlers

import (
	"net/http"
)

func handleEditMealPlan(w http.ResponseWriter, r *http.Request) {
	mpID, ok := getUint64Var(r, "mealplanid")
	if !ok {
		httpError(w, BadRequestError)
		return
	}
	
	renderTemplate(w, "edit-mp-form.html", mpID)
}

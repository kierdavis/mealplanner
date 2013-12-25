package mpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type JsonResponse struct {
	Error string `json:"error"`
	Success interface{} `json:"success"`
}

func HandleApiCall(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	
	err := r.ParseForm()
	if err != nil {
		response = JsonResponse{Error: "Could not parse request body."}
	} else {
		response = Dispatch(r.Form)
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not write JSON response: %s\n", err.Error())
	}
}

func Dispatch(params url.Values) (response JsonResponse) {
	switch params.Get("command") {
	case "fetch-meal-list":
		return fetchMealList(params)
	case "toggle-favourite":
		return toggleFavourite(params)
	case "delete-meal":
		return deleteMeal(params)
	}
	
	return JsonResponse{Error: "Invalid or missing command"}
}

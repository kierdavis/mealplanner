package mpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Type JsonResponse contains the response structure returned to the client.
// If the 'Error' field is nonempty, the response indicates an error has
// occurred, else the response is assumed to be a successful one.
type JsonResponse struct {
	Error   string      `json:"error"`   // The error message, in the event of an unsuccessful response.
	Success interface{} `json:"success"` // The response payload, in the event of a successful response.
}

// HandleApiCall handles an HTTP request for an API call. It obtains the form
// values, passes them through Dispatch and sends the resulting JSON response
// to the client.
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

// Dispatch looks at the "command" parameter and dispatches the request to the
// appropriate handler function.
func Dispatch(params url.Values) (response JsonResponse) {
	switch params.Get("command") {
	case "fetch-meal-list":
		return fetchMealList(params)
	case "toggle-favourite":
		return toggleFavourite(params)
	case "delete-meal":
		return deleteMeal(params)
	case "fetch-all-tags":
		return fetchAllTags(params)
	case "fetch-servings":
		return fetchServings(params)
	case "fetch-suggestions":
		return fetchSuggestions(params)
	case "update-serving":
		return updateServing(params)
	case "delete-serving":
		return deleteServing(params)
	}

	return JsonResponse{Error: "Invalid or missing command"}
}

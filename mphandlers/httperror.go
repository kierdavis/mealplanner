package mphandlers

import (
	"net/http"
)

// HTTPError holds related information about an HTTP status code used by the
// application.
type HTTPError struct {
	Status    int    // The HTTP status code.
	ShortDesc string // The associated "reason" message sent with the status code.
	LongDesc  string // A longer message displayed to the user on the HTML error page.
}

// BadRequestError represents an HTTP 400 Bad Request error.
var BadRequestError = &HTTPError{
	Status:    http.StatusBadRequest,
	ShortDesc: "Bad Request",
	LongDesc:  "We're sorry, there was an error when processing your request.",
}

// NotFoundError represents an HTTP 404 Not Found error.
var NotFoundError = &HTTPError{
	Status:    http.StatusNotFound,
	ShortDesc: "Not Found",
	LongDesc:  "We're sorry, the page you were looking for was not found on ther server.",
}

// InternalServerError represents an HTTP 500 Internal Server Error.
var InternalServerError = &HTTPError{
	Status:    http.StatusInternalServerError,
	ShortDesc: "Internal Server Error",
	LongDesc:  "We're sorry, the server encountered an unexpected error and was unable to complete the request.",
}

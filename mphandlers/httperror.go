package mphandlers

// Type HttpError holds related information about an HTTP status code used by
// the application.
type HttpError struct {
	Status    int    // The HTTP status code.
	ShortDesc string // The associated "reason" message sent with the status code.
	LongDesc  string // A longer message displayed to the user on the HTML error page.
}

var BadRequestError = &HttpError{
	Status:    400,
	ShortDesc: "Bad Request",
	LongDesc:  "We're sorry, there was an error when processing your request.",
}

var InternalServerError = &HttpError{
	Status:    500,
	ShortDesc: "Internal Server Error",
	LongDesc:  "We're sorry, the server encountered an unexpected error and was unable to complete the request.",
}

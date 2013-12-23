package mphandlers

type HttpError struct {
	Status int
	ShortDesc string
	LongDesc string
}

var InvalidMethodError = &HttpError{
	Status: 405,
	ShortDesc: "Invalid Method",
	LongDesc: "This page does not support the HTTP method that was used to request it.",
}

var InternalServerError = &HttpError{
	Status: 500,
	ShortDesc: "Internal Server Error",
	LongDesc: "We're sorry, the server encountered an unexpected error and was unable to complete the request.",
}

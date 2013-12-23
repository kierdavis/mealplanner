package mphandlers

type HttpError struct {
	Status int
	ShortDesc string
	LongDesc string
}

var BadRequestError = &HttpError{
	Status: 400,
	ShortDesc: "Bad Request",
	LongDesc: "We're sorry, there was an error when processing your request.",
}

var InternalServerError = &HttpError{
	Status: 500,
	ShortDesc: "Internal Server Error",
	LongDesc: "We're sorry, the server encountered an unexpected error and was unable to complete the request.",
}

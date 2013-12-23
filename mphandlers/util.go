package mphandlers

import (
	"fmt"
	"github.com/kierdavis/mealplanner/mptemplates"
	"net/http"
	"os"
	"runtime"
)

func httpError(w http.ResponseWriter, h *HttpError) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	w.WriteHeader(h.Status)
	
	err := mptemplates.Templates.ExecuteTemplate(w, "error.html", h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Internal error when rendering error.html: %s\n", err.Error())
	}
}

func serverError(w http.ResponseWriter, err error) {
	fmt.Fprintf(os.Stderr, "Internal error: %s\n", err.Error())
	
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(os.Stderr, "  at %s line %d\n", filename, lineno)
	}
	
	httpError(w, InternalServerError)
}

func invalidMethod(w http.ResponseWriter, methods string) {
	w.Header().Set("Allow", methods)
	httpError(w, InvalidMethodError)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	
	err := mptemplates.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		serverError(w, err)
	}
}

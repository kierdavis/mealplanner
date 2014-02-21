package mphandlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kierdavis/mealplanner/mpresources"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

// httpError sends an HTTP error code to the client followed by an HTML error
// page explaining the error.
func httpError(w http.ResponseWriter, h *HTTPError) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	w.WriteHeader(h.Status)

	err := mpresources.Templates.ExecuteTemplate(w, "error.html", h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Internal error when rendering error.html: %s\n", err.Error())
	}
}

// serverError logs 'err' to the console, then sends a 500 Internal Server Error
// response to the client.
func serverError(w http.ResponseWriter, err error) {
	fmt.Fprintf(os.Stderr, "Internal error: %s\n", err.Error())

	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(os.Stderr, "  at %s line %d\n", filename, lineno)
	}

	httpError(w, InternalServerError)
}

// redirect sends a redirection response to the client with the given status
// code.
func redirect(w http.ResponseWriter, status int, url string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	w.Header().Set("Location", url)
	w.WriteHeader(status)
	fmt.Fprintf(w, "Redirecting to %s...\n", url)
}

// renderTemplate renders the named template and returns the rendered HTML to
// the client.
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")

	err := mpresources.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		serverError(w, err)
	}
}

// getUint64Var gets a URL variable (a parameter embedded in the request URI)
// and parses it as an unsigned 64-bit integer.
func getUint64Var(r *http.Request, name string) (value uint64, ok bool) {
	vars := mux.Vars(r)
	str, ok := vars[name]
	if !ok {
		return 0, false
	}

	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, false
	}

	return value, true
}

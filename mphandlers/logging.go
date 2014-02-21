package mphandlers

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingHandler wraps an http.Handler, printing a message to standard output
// whenever a request is handled.
type LoggingHandler struct {
	Handler http.Handler
}

// ServeHTTP handles an HTTP request.
func (lh LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggingWriter{writer: w, status: 200}
	lh.Handler.ServeHTTP(lw, r)
	lh.logRequest(lw, r)
}

// logRequest logs a request to standard output.
func (lh LoggingHandler) logRequest(lw *loggingWriter, r *http.Request) {
	t := time.Now().Format("2006-01-02 15:04:05 -07:00")
	fmt.Printf("[%s] %s %s -> %d %s from %s\n", t, r.Method, r.URL.String(), lw.status, http.StatusText(lw.status), r.RemoteAddr)
}

// loggingWriter wraps an http.ResponseWriter, recording the status code sent
// to the client and the total number of bytes written.
type loggingWriter struct {
	writer       http.ResponseWriter
	status       int
	totalWritten int
}

// Header returns the header map that will be sent by WriteHeader.
func (lw *loggingWriter) Header() (h http.Header) {
	return lw.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (lw *loggingWriter) Write(buffer []byte) (n int, err error) {
	n, err = lw.writer.Write(buffer)
	lw.totalWritten += n
	return n, err
}

// WriteHeader sends an HTTP response header with status code.
func (lw *loggingWriter) WriteHeader(status int) {
	lw.status = status
	lw.writer.WriteHeader(status)
}

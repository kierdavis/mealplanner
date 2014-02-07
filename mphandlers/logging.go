package mphandlers

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingHandler struct {
	Handler http.Handler
}

func (lh LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggingWriter{writer: w}
	lh.Handler.ServeHTTP(lw, r)
	lh.logRequest(lw, r)
}

func (lh LoggingHandler) logRequest(lw *loggingWriter, r *http.Request) {
	t := time.Now().Format("2006-01-02 15:04:05 -07:00")
	fmt.Printf("[%s] %s %s -> %d %s from %s\n", t, r.Method, r.URL.String(), lw.status, http.StatusText(lw.status), r.RemoteAddr)
}

type loggingWriter struct {
	writer       http.ResponseWriter
	status       int
	totalWritten int
}

func (lw *loggingWriter) Header() (h http.Header) {
	return lw.writer.Header()
}

func (lw *loggingWriter) Write(buffer []byte) (n int, err error) {
	n, err = lw.writer.Write(buffer)
	lw.totalWritten += n
	return n, err
}

func (lw *loggingWriter) WriteHeader(status int) {
	lw.status = status
	lw.writer.WriteHeader(status)
}

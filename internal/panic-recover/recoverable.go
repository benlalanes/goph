package panic_recover

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Recoverable struct {
	buf *bytes.Buffer
	handler http.Handler
	header http.Header
	statusCode int
}

func (r *Recoverable) Header() http.Header {
	return r.header
}

func (r *Recoverable) Write(p []byte) (int, error) {
	return r.buf.Write(p)
}

func (r *Recoverable) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func NewRecoverable(handler http.Handler) *Recoverable {

	return &Recoverable{
		buf: &bytes.Buffer{},
		handler: handler,
		header: make(http.Header),
		statusCode: 200,
	}
}

func (r *Recoverable) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// If there is a panic for any reason, recover the panic, send 500
	// status code, and display error message.
	defer func() {
		if rec := recover(); rec != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "Something went wrong.")
		}

		r.buf.Reset()
	}()

	// Defer to wrapped handler.
	r.handler.ServeHTTP(r, req)

	// Write status code.
	w.WriteHeader(r.statusCode)

	// Copy headers.
	for k, v := range r.header {
		w.Header()[k] = v
	}

	// Copy response body.
	_, err := io.Copy(w, r.buf)
	if err != nil {
		panic(err)
	}

}
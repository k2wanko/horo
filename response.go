package horo

import (
	"bufio"
	"context"
	"net"
	"net/http"
)

type (
	// ResponseWriter is horo ResponseWriter
	ResponseWriter interface {
		http.ResponseWriter
		http.Hijacker
		http.Flusher
		http.CloseNotifier

		Status() int
		Size() int64
		Committed() bool
		Reset(http.ResponseWriter)
	}

	response struct {
		http.ResponseWriter
		code      int
		size      int64
		committed bool
	}
)

// Response returns ResponseWriter from context.
func Response(c context.Context) (w ResponseWriter) {
	if c := fromCtx(c); c != nil {
		w = c.w
	}
	return
}

func (r *response) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *response) Write(b []byte) (n int, err error) {
	if !r.committed {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	return
}

func (r *response) WriteHeader(code int) {
	if r.committed {
		return
	}
	r.code = code
	r.ResponseWriter.WriteHeader(code)
	r.committed = true
}

func (r *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *response) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

func (r *response) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (r *response) Status() int {
	return r.code
}

func (r *response) Size() int64 {
	return r.size
}

func (r *response) Committed() bool {
	return r.committed
}

func (r *response) Reset(w http.ResponseWriter) {
	r.ResponseWriter = w
	r.size = 0
	r.code = 0
	r.committed = false
}

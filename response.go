package horo

import (
	"bufio"
	"net"
	"net/http"

	"golang.org/x/net/context"
)

type (
	// ResponseWriter is horo ResponseWriter
	ResponseWriter struct {
		rw        http.ResponseWriter
		code      int
		size      int64
		committed bool
	}
)

// Response returns ResponseWriter from context.
func Response(c context.Context) (w *ResponseWriter) {
	if c := fromCtx(c); c != nil {
		w = c.w
	}
	return
}

func (r *ResponseWriter) Header() http.Header {
	return r.rw.Header()
}

func (r *ResponseWriter) Write(b []byte) (n int, err error) {
	if !r.committed {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.rw.Write(b)
	r.size += int64(n)
	return
}

func (r *ResponseWriter) WriteHeader(code int) {
	if r.committed {
		return
	}
	r.code = code
	r.rw.WriteHeader(code)
	r.committed = true
}

func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.rw.(http.Hijacker).Hijack()
}

func (r *ResponseWriter) Flush() {
	r.rw.(http.Flusher).Flush()
}

func (r *ResponseWriter) CloseNotify() <-chan bool {
	return r.rw.(http.CloseNotifier).CloseNotify()
}

func (r *ResponseWriter) Status() int {
	return r.code
}

func (r *ResponseWriter) Size() int64 {
	return r.size
}

func (r *ResponseWriter) Committed() bool {
	return r.committed
}

func (r *ResponseWriter) Reset(w http.ResponseWriter) {
	r.rw = w
	r.size = 0
	r.code = 0
	r.committed = false
}

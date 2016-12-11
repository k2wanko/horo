/*
Package horo is context friendly, Simple Web framework.

Standard Example:

    package main

    import (
        "net/http"

        "golang.org/x/net/context"

        "github.com/k2wanko/horo"
    )

    func Index(c context.Context) error {
        return horo.Text(c, http.StatusOK, "Hello World!")
    }

    func main() {
        h := horo.New()

        h.GET("/", Index)

        h.ListenAndServe(":8080")
    }

Google App Engine Example:

    package main

    import (
        "net/http"

        "golang.org/x/net/context"

        "github.com/k2wanko/horo"
    )

    func Index(c context.Context) error {
        return horo.Text(c, http.StatusOK, "Hello World!")
    }

    func init() {
        h := horo.New()

        h.GET("/", Index)

        http.Handle("/", h)
    }
*/
package horo

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

type (
	// Horo freamwork instance.
	Horo struct {
		ErrorHandler ErrorHandlerFunc

		router     *httprouter.Router
		middleware []MiddlewareFunc
	}

	// HandlerFunc is server HTTP requests.
	HandlerFunc func(context.Context) error

	// MiddlewareFunc is process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// ErrorHandlerFunc is error handling function.
	ErrorHandlerFunc func(context.Context, error)

	// HTTPError handling a request.
	HTTPError struct {
		Code    int
		Message string
	}
)

var (
	// ErrNotContext is thrown if the context does not have Horo context.
	ErrNotContext = errors.New("Not Context")

	// ErrInvalidRedirectCode is thrown if invalid redirect code.
	ErrInvalidRedirectCode = errors.New("invalid redirect status code")
)

// New is create Horo instance.
func New() (h *Horo) {
	h = &Horo{
		ErrorHandler: DefaultErrorHandler,
		router:       httprouter.New(),
	}
	return
}

// Use is add middleware.
func (h *Horo) Use(mw ...MiddlewareFunc) {
	h.middleware = append(h.middleware, mw...)
}

// GET registers a new GET handler
func (h *Horo) GET(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.GET(path, hf.hrHandle(h, mw...))
}

// POST registers a new POST handler
func (h *Horo) POST(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.POST(path, hf.hrHandle(h, mw...))
}

// PATCH registers a new PATCH handler
func (h *Horo) PATCH(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.PATCH(path, hf.hrHandle(h, mw...))
}

// PUT registers a new PUT handler
func (h *Horo) PUT(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.PUT(path, hf.hrHandle(h, mw...))
}

// OPTIONS registers a new OPTIONS handler
func (h *Horo) OPTIONS(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.OPTIONS(path, hf.hrHandle(h, mw...))
}

// DELETE registers a new DELETE handler
func (h *Horo) DELETE(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.DELETE(path, hf.hrHandle(h, mw...))
}

// HEAD registers a new HEAD handler
func (h *Horo) HEAD(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.HEAD(path, hf.hrHandle(h, mw...))
}

func (hf HandlerFunc) hrHandle(h *Horo, mw ...MiddlewareFunc) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w := &response{ResponseWriter: rw}
		var c context.Context = &horoCtx{
			Context: context.Background(),
			w:       w,
			r:       r,
			ps:      ps,
		}

		c, cancel := context.WithCancel(c)

		hf := func(c context.Context) error {
			mw := append(h.middleware, mw...)
			for i := len(mw) - 1; i >= 0; i-- {
				hf = mw[i](hf)
			}

			return hf(c)
		}

		if err := hf(c); err != nil {
			h.ErrorHandler(c, err)
		}

		cancel()
	}
}

// DefaultErrorHandler invoke HTTP Error Handler
func DefaultErrorHandler(c context.Context, err error) {
	code := 500
	msg := http.StatusText(code)
	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	Text(c, code, msg)
}

// ServeHTTP implements http.Handler interface.
func (h *Horo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// ListenAndServe is Start HTTP Server
func (h *Horo) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (e *HTTPError) Error() string {
	return e.Message
}

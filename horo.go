/*
Package horo is context friendly, Simple Web framework.

Basic Example:
    package main

    import (
        "net/http"

        "golang.org/x/net/context"

        "github.com/k2wanko/horo"
        "github.com/k2wanko/horo/middleware"
    )

    func Index(c context.Context) error {
        return horo.Text(c, http.StatusOK, "Hello World!")
    }

    func main() {
        h := horo.New()
        h.Use(middleware.Logger())

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
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/k2wanko/horo/log"
	"golang.org/x/net/context"
)

type (
	// Horo freamwork instance.
	Horo struct {
		ErrorHandler ErrorHandlerFunc
		Logger       log.Logger

		router     *httprouter.Router
		middleware []MiddlewareFunc
		pool       sync.Pool
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
		Logger:       log.DefaultLogger,
		router:       httprouter.New(),
	}

	h.pool.New = func() interface{} {
		return &horoCtx{
			Context: context.Background(),
			w:       &response{},
		}
	}

	return
}

// Use is add middleware.
func (h *Horo) Use(mw ...MiddlewareFunc) {
	h.middleware = append(h.middleware, mw...)
}

// GET registers a new GET handler
func (h *Horo) GET(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.GET(path, h.handle(hf, mw...))
}

// POST registers a new POST handler
func (h *Horo) POST(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.POST(path, h.handle(hf, mw...))
}

// PATCH registers a new PATCH handler
func (h *Horo) PATCH(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.PATCH(path, h.handle(hf, mw...))
}

// PUT registers a new PUT handler
func (h *Horo) PUT(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.PUT(path, h.handle(hf, mw...))
}

// OPTIONS registers a new OPTIONS handler
func (h *Horo) OPTIONS(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.OPTIONS(path, h.handle(hf, mw...))
}

// DELETE registers a new DELETE handler
func (h *Horo) DELETE(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.DELETE(path, h.handle(hf, mw...))
}

// HEAD registers a new HEAD handler
func (h *Horo) HEAD(path string, hf HandlerFunc, mw ...MiddlewareFunc) {
	h.router.HEAD(path, h.handle(hf, mw...))
}

func (h *Horo) handle(hf HandlerFunc, mwf ...MiddlewareFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		hc := h.pool.Get().(*horoCtx)
		hc.Reset(w, r, ps)

		c, cancel := context.WithCancel(hc)
		c = log.WithContext(c, h.Logger)

		hwl := len(h.middleware)
		mw := make([]MiddlewareFunc, hwl+len(mwf))
		for i := 0; i < cap(mw); i++ {
			if i < hwl {
				mw[i] = h.middleware[i]
			} else {
				mw[i] = mwf[i-hwl]
			}
		}

		hf := func(c context.Context) error {
			for i := len(mw) - 1; i >= 0; i-- {
				hf = mw[i](hf)
			}

			return hf(c)
		}

		if err := hf(c); err != nil {
			h.ErrorHandler(c, err)
		}

		cancel()

		h.pool.Put(hc)
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

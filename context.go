package horo

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type (
	// RequestIDGenerator is generate requestID method interface.
	RequestIDGenerator interface {
		RequestID(c context.Context) string
	}

	horoCtx struct {
		context.Context

		h     *Horo
		w     *ResponseWriter
		r     *http.Request
		ps    httprouter.Params
		reqID string
	}

	ctxkey struct {
		name string
	}

	reqGen struct{}
)

var (
	// DefaultRequestIDGenerator is application request id generator.
	DefaultRequestIDGenerator RequestIDGenerator = &reqGen{}

	ctxKey = &ctxkey{"horo ctx"}
)

func fromCtx(c context.Context) (ctx *horoCtx) {
	ctx, _ = c.Value(ctxKey).(*horoCtx)
	return
}

func (c *horoCtx) Value(key interface{}) interface{} {
	if key == ctxKey {
		return c
	}
	return c.Context.Value(key)
}

func (c *horoCtx) Reset(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.w.Reset(rw)
	c.r = r
	c.ps = ps
	c.reqID = ""
}

// Param returns url param.
func Param(c context.Context, name string) (v string) {
	if c := fromCtx(c); c != nil {
		if c.ps != nil {
			v = c.ps.ByName(name)
		}
	}
	return
}

// RequestID returns request id from context.
func RequestID(ctx context.Context) (id string) {
	if c := fromCtx(ctx); c != nil {
		id = c.reqID
		if id != "" {
			return
		}

		id = requestID(c.r)
		if g := c.h.RequestIDGenerator; id == "" && g != nil {
			id = g.RequestID(ctx)
		} else {
			id = DefaultRequestIDGenerator.RequestID(ctx)
		}
		c.reqID = id
	}
	return
}

func (*reqGen) RequestID(c context.Context) string {
	return uuid.NewV4().String()
}

// NoContent send no body.
func NoContent(c context.Context, code int) error {
	if c := fromCtx(c); c != nil {
		c.w.WriteHeader(code)
		return nil
	}
	return ErrNotContext
}

// Text send a Text response.
func Text(c context.Context, code int, s string) (err error) {
	if c := fromCtx(c); c != nil {
		c.w.Header().Add("Content-Type", "text/plain")
		c.w.WriteHeader(code)
		_, err = c.w.Write([]byte(s))
		return
	}
	return ErrNotContext
}

// HTML send a HTML response.
func HTML(c context.Context, code int, html string) (err error) {
	if c := fromCtx(c); c != nil {
		c.w.Header().Add("Content-Type", "text/html")
		c.w.WriteHeader(code)
		_, err = c.w.Write([]byte(html))
		return
	}
	return ErrNotContext
}

// JSON send a JSON response.
func JSON(c context.Context, code int, i interface{}) (err error) {
	if c := fromCtx(c); c != nil {
		c.w.Header().Add("Content-Type", "application/json")
		c.w.WriteHeader(code)
		var b []byte
		b, err = json.Marshal(i)
		if err != nil {
			return
		}
		_, err = c.w.Write(b)
		return
	}
	return ErrNotContext
}

// Redirect redirect the request status code.
func Redirect(c context.Context, code int, url string) error {
	if c := fromCtx(c); c != nil {
		if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
			return nil
		}
		c.w.Header().Set("Location", url)
		c.w.WriteHeader(code)
		return nil
	}
	return ErrNotContext
}

// Request returns *http.Request from context.
func Request(c context.Context) (r *http.Request) {
	if c := fromCtx(c); c != nil {
		r = c.r
	}
	return
}

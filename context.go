package horo

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

type (
	horoCtx struct {
		context.Context

		w  ResponseWriter
		r  *http.Request
		ps httprouter.Params
	}

	ctxkey struct {
		name string
	}
)

var (
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

// Param returns url param.
func Param(c context.Context, name string) string {
	if c := fromCtx(c); c != nil {
		return c.ps.ByName(name)
	}
	return ""
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

// Response returns http.ResponseWriter from context.
func Response(c context.Context) (w ResponseWriter) {
	if c := fromCtx(c); c != nil {
		w = c.w
	}
	return
}

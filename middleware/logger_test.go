//+build !appengine

package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/k2wanko/horo"
	"github.com/k2wanko/horo/log"
)

func TestLogger(t *testing.T) {
	h := horo.New()

	out := new(bytes.Buffer)

	h.Logger = log.New(log.Out(out))

	h.Use(Logger())

	h.GET("/", func(c context.Context) error {
		return horo.Text(c, 200, "Test")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	t.Logf("\n%s", out)
	if out := out.String(); out == "" {
		t.Error("Log empty")
	}
}

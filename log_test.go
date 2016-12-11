//+build !appengine

package horo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/k2wanko/horo/log"
)

func TestLog(t *testing.T) {
	out := new(bytes.Buffer)
	testLogger := log.New(log.Out(out))

	h := New()

	// Set Logger
	h.Use(func(next HandlerFunc) HandlerFunc {
		return func(c context.Context) error {
			c = log.WithContext(c, testLogger)
			return next(c)
		}
	})

	h.GET("/", func(c context.Context) error {
		l := log.FromContext(c)
		l.Infof(c, "Test")
		return nil
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if out, want := out.String(), "[INFO] Test\n"; out != want {
		t.Errorf("out = %s; want = %s", out, want)
	}
}

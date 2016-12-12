package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k2wanko/horo"
	"golang.org/x/net/context"
)

func TestRecover(t *testing.T) {
	h := horo.New()
	call := false
	h.Use(func(next horo.HandlerFunc) horo.HandlerFunc {
		return func(c context.Context) (err error) {
			err = next(c)
			call = true
			if err.Error() != "Test" {
				t.Errorf("err = %s; want = %s", err.Error(), "Test")
			}
			return
		}
	})
	h.Use(Recover())

	h.GET("/", func(c context.Context) error {
		panic("Test")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 500; code != want {
		t.Errorf("code = %v; want = %v", code, want)
	}

	if out, want := w.Body.String(), "Internal Server Error"; out != want {
		t.Errorf("out = %s; want = %s", out, want)
	}

	if !call {
		t.Error("Not call Top level middleware.")
	}
}

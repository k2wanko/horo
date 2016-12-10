package horo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

func TestParam(t *testing.T) {
	h := New()
	h.GET("/:user", func(c context.Context) error {
		return Text(c, 200, Param(c, "user"))
	})

	r, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 200; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), "me"; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

func TestJSON(t *testing.T) {
	h := New()
	h.GET("/", func(c context.Context) error {
		return JSON(c, 200, map[string]string{"user": "k2wanko"})
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 200; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), `{"user":"k2wanko"}`; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

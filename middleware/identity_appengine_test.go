//+build appengine

package middleware

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/k2wanko/horo"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func TestAppContext(t *testing.T) {
	i, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer i.Close()

	h := horo.New()
	h.Use(AppContext())
	h.GET("/", func(c context.Context) error {
		return horo.Text(c, 200, fmt.Sprintf("Hello, %s", appengine.AppID(c)))
	})
	r, _ := i.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 200; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), "Hello, testapp"; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

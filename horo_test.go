package horo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"golang.org/x/net/context"
)

func TestSimpleHandle(t *testing.T) {
	h := New()
	h.GET("/", func(c context.Context) error {
		return Text(c, 200, "ok")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 200; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), "ok"; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

func TestNotFound(t *testing.T) {
	h := New()
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 404; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), "Not Found"; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

func TestErrorHandler(t *testing.T) {
	h := New()
	h.GET("/", func(c context.Context) error {
		return errors.New("Test")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 500; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), http.StatusText(500); body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

func TestHTTPError(t *testing.T) {
	h := New()
	h.GET("/", func(c context.Context) error {
		return &HTTPError{Code: 501, Message: "Test"}
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if code, want := w.Code, 501; code != want {
		t.Errorf("w.Code = %v; want %v", code, want)
	}

	if body, want := w.Body.String(), "Test"; body != want {
		t.Errorf("w.Body = %v; want %v", body, want)
	}
}

func TestMiddleware(t *testing.T) {
	var res []int
	h := New()
	h.Use(func(next HandlerFunc) HandlerFunc {
		return func(c context.Context) (err error) {
			t.Log("Top level middleware: before")
			res = append(res, 0)
			err = next(c)
			t.Log("Top level middleware: after")
			res = append(res, 4)
			return
		}
	})
	h.GET("/", func(c context.Context) error {
		t.Log("Handler")
		res = append(res, 2)
		return nil
	}, func(next HandlerFunc) HandlerFunc {
		return func(c context.Context) (err error) {
			t.Log("Handler middleware: before")
			res = append(res, 1)
			err = next(c)
			t.Log("Handler middleware: after")
			res = append(res, 3)
			return
		}
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if diff := pretty.Compare(res, []int{0, 1, 2, 3, 4}); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

func TestMiddlewareNoHandler(t *testing.T) {
	var res []int
	h := New()
	h.Use(func(next HandlerFunc) HandlerFunc {
		return func(c context.Context) (err error) {
			t.Log("Top level middleware: before")
			res = append(res, 0)
			err = next(c)
			t.Log("Top level middleware: after")
			res = append(res, 1)
			return
		}
	})

	h.GET("/", func(c context.Context) error {
		res = append(res, -1)
		return Text(c, 200, "Test")
	})

	r, _ := http.NewRequest("GET", "/404", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if diff := pretty.Compare(res, []int{0, 1}); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}

	if code, want := w.Code, 404; code != want {
		t.Errorf("code = %v; want = %v", code, want)
	}

	if out, want := w.Body.String(), "Not Found"; out != want {
		t.Errorf("out = %s; want = %s", out, want)
	}

	// reset
	res = nil
	r, _ = http.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if diff := pretty.Compare(res, []int{0, 1}); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}

	if code, want := w.Code, 405; code != want {
		t.Errorf("code = %v; want = %v", code, want)
	}

	if out, want := w.Body.String(), "Method Not Allowed"; out != want {
		t.Errorf("out = %s; want = %s", out, want)
	}
}

type mockResponseWriter struct {
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteHeader(code int) {

}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := mockResponseWriter{}
	u := r.URL
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(&w, r)

		// clear caches
		r.Form = nil
		r.PostForm = nil
		r.MultipartForm = nil
	}
}

var testHandler http.Handler

func init() {
	h := New()
	h.GET("/", func(c context.Context) error {
		return Text(c, 200, "Hello, World")
	})
	h.GET("/:name", func(c context.Context) error {
		return Text(c, 200, fmt.Sprintf("Hello, %s", Param(c, "name")))
	})
	testHandler = h
}

func BenchmarkRequest(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)
	benchRequest(b, testHandler, req)
}

func BenchmarkRequestParam(b *testing.B) {
	req, _ := http.NewRequest("GET", "/gopher", nil)
	benchRequest(b, testHandler, req)
}

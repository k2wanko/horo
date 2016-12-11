package horo_test

import (
	"net/http"

	"github.com/k2wanko/horo"
	"github.com/k2wanko/horo/middleware"
	"golang.org/x/net/context"
)

func Example_basic() {
	h := horo.New()
	h.Use(middleware.Logger())
	h.GET("/", func(c context.Context) error {
		return horo.Text(c, 200, "Hello, World")
	})
	h.ListenAndServe(":8080")
}

func Example_appengine() {
	// Write in func init()
	h := horo.New()
	h.GET("/", func(c context.Context) error {
		return horo.Text(c, 200, "Hello, World")
	})
	http.Handle("/", h)
}

//+build appengine

package middleware

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/k2wanko/horo"
)

// AppContext is middleware that sets the context of App Engine.
func AppContext() horo.MiddlewareFunc {
	return func(next horo.HandlerFunc) horo.HandlerFunc {
		return func(c context.Context) error {
			return next(appengine.WithContext(c, horo.Request(c)))
		}
	}
}

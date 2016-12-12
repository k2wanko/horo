//+build !appengine

package middleware

import (
	"github.com/k2wanko/horo"
)

// AppContext is middleware that sets the context of App Engine.
// It does not work on platforms other than appengine.
func AppContext() horo.MiddlewareFunc {
	return func(next horo.HandlerFunc) horo.HandlerFunc {
		return next
	}
}

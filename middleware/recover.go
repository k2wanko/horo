package middleware

import (
	"fmt"

	"github.com/k2wanko/horo"
	"golang.org/x/net/context"
)

// Recover returns a middleware which recovers from panics
func Recover() horo.MiddlewareFunc {
	return func(next horo.HandlerFunc) horo.HandlerFunc {
		return func(c context.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
				}
			}()

			return next(c)
		}
	}
}

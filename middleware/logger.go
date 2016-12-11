package middleware

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/k2wanko/horo"
	"github.com/k2wanko/horo/log"
)

type (
	// RequestLogger is request base loggin.
	RequestLogger struct {
		log.Logger
	}
)

// Logger is horo Logger middleware.
func Logger() horo.MiddlewareFunc {
	return func(next horo.HandlerFunc) horo.HandlerFunc {
		return func(c context.Context) (err error) {
			start := time.Now()
			r := horo.Request(c)
			path := r.URL.Path

			l := &RequestLogger{
				Logger: log.FromContext(c),
			}
			c = log.WithContext(c, l)

			// process
			err = next(c)

			end := time.Now()
			latency := end.Sub(start)
			ip := r.RemoteAddr
			method := r.Method
			code := horo.Response(c).Status()
			l.Infof(c, "%v %3d %v %s %s %s",
				end.Format("2006/01/02 - 15:04:05"),
				code,
				latency,
				ip,
				method,
				path,
			)

			return
		}
	}
}

// Debugf implements log.Logger#Debugf
func (rl *RequestLogger) Debugf(c context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	rl.Logger.Debugf(c, "%s", msg)
}

// Infof implements log.Logger#Infof
func (rl *RequestLogger) Infof(c context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	rl.Logger.Infof(c, "%s", msg)
}

// Warnf implements log.Logger#Warnf
func (rl *RequestLogger) Warnf(c context.Context, format string, args ...interface{}) {

}

// Errorf implements log.Logger#Errorf
func (rl *RequestLogger) Errorf(c context.Context, format string, args ...interface{}) {

}

// Fatalf implements log.Logger#Fatalf
func (rl *RequestLogger) Fatalf(c context.Context, format string, args ...interface{}) {

}

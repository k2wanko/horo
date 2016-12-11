//+build appengine

package log

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

func (l *logger) Debugf(c context.Context, format string, args ...interface{}) {
	log.Debugf(c, format, args...)
}

func (l *logger) Infof(c context.Context, format string, args ...interface{}) {
	log.Infof(c, format, args...)
}

func (l *logger) Warnf(c context.Context, format string, args ...interface{}) {
	log.Warningf(c, format, args...)
}

func (l *logger) Errorf(c context.Context, format string, args ...interface{}) {
	log.Errorf(c, format, args...)
}

func (l *logger) Fatalf(c context.Context, format string, args ...interface{}) {
	log.Criticalf(c, format, args...)
}

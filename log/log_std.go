//+build !appengine

package log

import (
	"fmt"
	"os"

	"golang.org/x/net/context"
)

func (l *logger) Debugf(c context.Context, format string, args ...interface{}) {
	l.Write(c, DEBUG, format, args...)
}

func (l *logger) Infof(c context.Context, format string, args ...interface{}) {
	l.Write(c, INFO, format, args...)
}

func (l *logger) Warnf(c context.Context, format string, args ...interface{}) {
	l.Write(c, WARN, format, args...)
}

func (l *logger) Errorf(c context.Context, format string, args ...interface{}) {
	l.Write(c, ERROR, format, args...)
}

func (l *logger) Fatalf(c context.Context, format string, args ...interface{}) {
	l.Write(c, FATAL, format, args...)
}

func (l *logger) Write(c context.Context, lvl Level, format string, args ...interface{}) {
	if l.out == nil {
		l.out = os.Stdout
	}

	if l.err == nil {
		l.err = os.Stderr
	}

	w := l.out

	switch lvl {
	case WARN:
		fallthrough
	case ERROR:
		fallthrough
	case FATAL:
		w = l.err
	}

	fmt.Fprintf(w, "[%s] %s\n", lvl, fmt.Sprintf(format, args...))
}

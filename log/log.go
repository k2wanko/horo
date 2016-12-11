package log

import (
	"io"

	"golang.org/x/net/context"
)

type (
	// Logger is application logger.
	Logger interface {
		Debugf(c context.Context, format string, args ...interface{})
		Infof(c context.Context, format string, args ...interface{})
		Warnf(c context.Context, format string, args ...interface{})
		Errorf(c context.Context, format string, args ...interface{})
		Fatalf(c context.Context, format string, args ...interface{})
	}

	// Level is log level.
	Level int

	// Option is New option.
	Option func(*opts)

	options []Option

	opts struct {
		out, err io.Writer
	}

	logger struct {
		out, err io.Writer
	}

	ctxkey struct {
		name string
	}
)

const (
	// DEBUG is Log Level
	DEBUG Level = iota

	// INFO is Log Level
	INFO

	// WARN is Log Level
	WARN

	// ERROR is Log Level
	ERROR

	// FATAL is Log Level
	FATAL
)

var (
	// DefaultLogger is default logger.
	DefaultLogger Logger

	logText = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}

	loggerContextKey = ctxkey{"logger"}
)

func init() {
	DefaultLogger = New()
}

// String returns Level text.
func (l Level) String() string {
	s, ok := logText[l]
	if ok {
		return s
	}
	return ""
}

func (o options) Option() *opts {
	opts := new(opts)
	for _, o := range o {
		o(opts)
	}
	return opts
}

// Out set logger out.
func Out(w io.Writer) Option {
	return func(o *opts) {
		o.out = w
	}
}

// ErrOut set logger error out.
func ErrOut(w io.Writer) Option {
	return func(o *opts) {
		o.err = w
	}
}

// New retunrts Logger
func New(opt ...Option) Logger {
	o := options(opt).Option()
	return &logger{
		out: o.out,
		err: o.err,
	}
}

// WithContext is set Logger.
func WithContext(c context.Context, l Logger) context.Context {
	return context.WithValue(c, loggerContextKey, l)
}

// FromContext returns Logger from context.
func FromContext(c context.Context) Logger {
	if l, ok := c.Value(loggerContextKey).(Logger); ok {
		return l
	}
	return DefaultLogger
}

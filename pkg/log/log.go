package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	_ StdLogger = &Logger{}
)

type StdLogger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithLogger() *logrus.Logger
	WithError(err error) *logrus.Entry
}

// Level represents a log level.
type Level int32

const (
	// FatalLevel is used for undesired and unexpected events that
	// the program cannot recover from.
	FatalLevel Level = iota

	// ErrorLevel is used for undesired and unexpected events that
	// the program can recover from.
	ErrorLevel

	// WarnLevel is used for undesired but relatively expected events,
	// which may indicate a problem.
	WarnLevel

	// InfoLevel is used for general informational log messages.
	InfoLevel

	// DebugLevel is the lowest level of logging.
	// Debug logs are intended for debugging and development purposes.
	DebugLevel
)

// String is part of the fmt.Stringer interface.
//
// Used for testing and debugging purposes.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l Level) ToLogrusLevel() (logrus.Level, error) {
	switch l {
	case DebugLevel:
		return logrus.DebugLevel, nil
	case InfoLevel:
		return logrus.InfoLevel, nil
	case WarnLevel:
		return logrus.WarnLevel, nil
	case ErrorLevel:
		return logrus.ErrorLevel, nil
	case FatalLevel:
		return logrus.FatalLevel, nil
	default:
		return 0, fmt.Errorf("not a valid log Level: %q", l)
	}
}

// NewLogger creates and returns a new instance of Logger.
// Log level is set to DebugLevel by default.
func NewLogger(out io.Writer) *Logger {
	log := &logrus.Logger{
		Out: out,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
		Level: logrus.DebugLevel,
	}

	log.SetReportCaller(false)

	return &Logger{
		logger: log,
		entry:  logrus.NewEntry(log),
	}
}

// Logger logs message to io.Writer at various log levels.
type Logger struct {
	logger *logrus.Logger

	entry *logrus.Entry
}

func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.entry.WithError(err)
}

func (l *Logger) WithLogger() *logrus.Logger {
	return l.logger
}

// SetLevel sets the logger level.
// It panics if v is less than DebugLevel or greater than FatalLevel.
func (l *Logger) SetLevel(v Level) {
	lvl, err := v.ToLogrusLevel()
	if err != nil {
		panic(err)
	}

	l.logger.SetLevel(lvl)
}

// WithField sets logger fields
func (l *Logger) SetPrefix(value interface{}) {
	l.entry = l.entry.WithField("source", value)
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}

package log

import (
	"context"
	"io"
	"net"
	"os"
	"runtime"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/env"
)

const (
	logFieldRequestID   = "request_id"
	logFieldJobID       = "job_id"
	logFieldExecutionID = "execution_id"
	logFieldAppVersion  = "app_version"
	logFieldTeam        = "team_name"
	logFieldCapability  = "capability"
	logFieldTier        = "tiering"
	logFieldTrafficType = "production_traffic"
	envTest             = "test"
)

var (
	// DefaultLogger is the logger the eventbus uses to log, overwrite
	// this with SetLogger if you want to see detailed logs of what the package is doing.
	DefaultLogger Logger = NewLogger(io.Discard, "", "version", "none", "none", "tier", false)

	stdFields = logrus.Fields{
		"component":         "service",
		logFieldAppVersion:  "unknown",
		logFieldRequestID:   "none",
		logFieldJobID:       "none",
		logFieldExecutionID: "none",
		logFieldTeam:        "none",
		logFieldCapability:  "none",
		logFieldTier:        0,
		logFieldTrafficType: "none",
	}
)

// Logger encapsulates Flyt compliant logging as an interface
type Logger interface {
	// Writer returns a pipe that can be used in the standard logger when
	// go internal require a standard logger
	Writer() *io.PipeWriter

	// WithContext extracts context scoped values into the logger
	WithContext(context.Context) Logger

	Dup() *logrus.Entry

	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type stdLogger struct {
	*logrus.Entry
}

// Bootstrap creates a logger writing to the appropriate place for the given environment
// In "test" env we output to os.Stderr, otherwise a udp connection to localhost:9999 will
// be used to emit all logs
func Bootstrap(env, app, team, capability, tier string, traffic bool) Logger {
	l := NewLogger(os.Stderr, app, "unknown", team, capability, tier, traffic)
	switch env {
	case envTest:
		return l
	default:
		conn, err := net.Dial("udp", ":9999")
		if err != nil {
			l.Warningf("Unable to connect to udp for logging: %s", err)
			return l
		}
		// Attempt to close the connection when the logger is garbage collected
		runtime.SetFinalizer(conn, func(conn net.Conn) {
			conn.Close()
		})
		return NewLogger(conn, app, "unknown", team, capability, tier, traffic)
	}
}

func BootstrapWithNetworkAndAddress(env, network, address, app, version, team, capability, tier string, traffic bool) Logger {
	l := NewLogger(os.Stderr, app, version, team, capability, tier, traffic)
	switch env {
	case envTest:
		return l
	default:
		conn, err := net.Dial(network, address)
		if err != nil {
			l.Warningf("Unable to connect to for logging: %s. Network %s, address %s.", err, network, address)

			conn, err = net.Dial("udp", ":9999")
			if err != nil {
				l.Warningf("Unable to connect to localhost:9999 udp for logging: %s", err)
				return l
			}
		}
		// Attempt to close the connection when the logger is garbage collected
		runtime.SetFinalizer(conn, func(conn net.Conn) {
			conn.Close()
		})
		return NewLogger(conn, app, version, team, capability, tier, traffic)
	}
}

// NewLogger returns a flyt compliant logger outputting on the given writer
func NewLogger(out io.Writer, app, version, team, capability, tier string, traffic bool) Logger {
	return &stdLogger{
		Entry: newLogrusEntry(out).WithFields(logrus.Fields{
			"app":               app,
			logFieldAppVersion:  version,
			logFieldTeam:        team,
			logFieldCapability:  capability,
			logFieldTier:        tier,
			logFieldTrafficType: traffic,
		}),
	}
}

// SetLogger will set the default logger which will be used by the log package.
func SetLogger(out io.Writer, app, version, team, capability, tier string, traffic bool) {
	DefaultLogger = &stdLogger{
		Entry: newLogrusEntry(out).WithFields(logrus.Fields{
			"app":               app,
			logFieldAppVersion:  version,
			logFieldTeam:        team,
			logFieldCapability:  capability,
			logFieldTier:        tier,
			logFieldTrafficType: traffic,
		}),
	}
}

// StdLogger returns a flyt compliant logger for StdOut, without the need for parameters
func StdLogger() Logger {
	return &stdLogger{
		Entry: newLogrusEntry(os.Stdout),
	}
}

// SetStdLogger will set the default logger to the StdLogger which will be used by the log package.
func SetStdLogger() {
	DefaultLogger = StdLogger()
}

func SetSentryHook(client *sentry.Client) {
	l := stdLogger{
		Entry: DefaultLogger.Dup(),
	}

	l.Entry.Logger.AddHook(NewSentryHook(client))

	DefaultLogger = &l
}

func (l *stdLogger) WithContext(ctx context.Context) Logger {
	e := l.Entry
	id, ok := flytcontext.RequestID(ctx)
	if ok {
		e = e.WithField(logFieldRequestID, id)
	}

	executionID, ok := flytcontext.ExecutionID(ctx)
	if ok {
		e = e.WithField(logFieldExecutionID, executionID)
	}

	jobID, ok := flytcontext.JobID(ctx)
	if ok {
		e = e.WithField(logFieldJobID, jobID)
	}

	return &stdLogger{
		Entry: e,
	}
}

func newLogrusEntry(out io.Writer) *logrus.Entry {
	logger := &logrus.Logger{
		Out:       out,
		Formatter: resolveLogFormatter(),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.SetReportCaller(true)
	return logger.WithFields(stdFields)
}

func resolveLogFormatter() logrus.Formatter {
	if env.IsLocal() {
		return &localFormatter{}
	}

	return flytJSONFormatter()
}

func flytJSONFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "message",
			logrus.FieldKeyFile: "app_line",
		},
	}
}

// Logger functions that wrap the default logger.

func Writer() *io.PipeWriter {
	return DefaultLogger.Writer()
}

func WithContext(ctx context.Context) Logger {
	return DefaultLogger.WithContext(ctx)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Print(args ...interface{}) {
	DefaultLogger.Print(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Warning(args ...interface{}) {
	DefaultLogger.Warning(args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	DefaultLogger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	DefaultLogger.Panicf(format, args...)
}

func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}

func Println(args ...interface{}) {
	DefaultLogger.Println(args...)
}

func Warnln(args ...interface{}) {
	DefaultLogger.Warnln(args...)
}

func Warningln(args ...interface{}) {
	DefaultLogger.Warningln(args...)
}

func Errorln(args ...interface{}) {
	DefaultLogger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	DefaultLogger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	DefaultLogger.Panicln(args...)
}

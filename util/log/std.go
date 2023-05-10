package log

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

var m = map[LogLevel]log.Level{
	Debug: log.LevelDebug,
	Error: log.LevelError,
	Fatal: log.LevelFatal,
	Info:  log.LevelInfo,
	Warn:  log.LevelWarn,
}

type StdLogger struct {
	Logger log.Logger
	level  LogLevel
}

func WrapLogger(log log.Logger, level LogLevel) Logger {
	return &StdLogger{
		Logger: log,
		level:  level,
	}
}

func (l *StdLogger) Level() LogLevel {
	return l.level
}

func (l *StdLogger) Log(level LogLevel, args ...interface{}) {
	if level < l.level {
		return
	}
	l.Logger.Log(getLevel(level), args)
}

func getLevel(wlevel LogLevel) log.Level {
	level, ok := m[wlevel]
	if !ok {
		return log.LevelDebug
	}
	return level
}

func (l *StdLogger) Info(args ...interface{}) {
	l.Log(Info, args)
}

func (l *StdLogger) Debug(args ...interface{}) {
	l.Log(Debug, args)
}

func (l *StdLogger) Error(args ...interface{}) {
	l.Log(Error, args)
}

func (l *StdLogger) Fatal(args ...interface{}) {
	l.Log(Fatal, args)
}

func (l *StdLogger) Warn(args ...interface{}) {
	l.Log(Warn, args)
}

func (l *StdLogger) Trace(args ...interface{}) {
	l.Log(Debug, args)
}

func (l *StdLogger) All(args ...interface{}) {
	l.Log(Debug, args)
}

func NewStdLogger(level LogLevel, serviceId string, serviceName string, serviceVersion string) Logger {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serviceId,
		"service.name", serviceName,
		"service.version", serviceVersion,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	return WrapLogger(logger, level)
}

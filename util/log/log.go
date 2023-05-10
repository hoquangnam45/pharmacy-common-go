package log

type LogLevel int

const (
	All LogLevel = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Trace
)

type Logger interface {
	Level() LogLevel
	Log(level LogLevel, args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	All(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Trace(args ...interface{})
}

package core

type LogLevel int

const (
	// !nashtsai! following level also match syslog.Priority value
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_WARNING
	LOG_ERR
	LOG_OFF
	LOG_UNKNOWN
)

// logger interface
type ILogger interface {
	Debug(v ...interface{}) (err error)
	Debugf(format string, v ...interface{}) (err error)
	Err(v ...interface{}) (err error)
	Errf(format string, v ...interface{}) (err error)
	Info(v ...interface{}) (err error)
	Infof(format string, v ...interface{}) (err error)
	Warning(v ...interface{}) (err error)
	Warningf(format string, v ...interface{}) (err error)

	Level() LogLevel
	SetLevel(l LogLevel) (err error)

	ShowSQL(show ...bool)
	IsShowSQL() bool
}

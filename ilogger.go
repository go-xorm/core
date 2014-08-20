package core

// logger interface, log/syslog conform with this interface
type ILogger interface {
	Debug(m string) (err error)
	Err(m string) (err error)
	Info(m string) (err error)
	Warning(m string) (err error)
}

package models

// Ilogger interface
type ILogger interface {
	Info(message string, v ...interface{})
	Error(message string, v ...interface{})
	Warn(message string, v ...interface{})
	Debug(message string, v ...interface{})
}

func NewLogger(l ILogger) ILogger {
	return l
}

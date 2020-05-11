package gocon

import "github.com/sirupsen/logrus"

const (
	FATAL = "fatal"
	ERROR = "error"
	WARN  = "warning"
	INFO  = "info"
	DEBUG = "debug"
	TRACE = "trace"
)

type ILogger interface {
	Debug(message interface{})
	Info(message interface{})
	Warn(message interface{})
	Error(message interface{})
	Fatal(message interface{})
	Trace(message interface{})
	Log(level string, message interface{})
}

//Default logger
type Instance struct {
	logger *logrus.Logger
}

func (l *Instance) Debug(message interface{}) {
	l.Log(DEBUG, message)
}

func (l *Instance) Info(message interface{}) {
	l.Log(INFO, message)
}

func (l *Instance) Warn(message interface{}) {
	l.Log(WARN, message)
}

func (l *Instance) Error(message interface{}) {
	l.Log(ERROR, message)
}

func (l *Instance) Fatal(message interface{}) {
	l.Log(FATAL, message)
}

func (l *Instance) Trace(message interface{}) {
	l.Log(TRACE, message)
}

func (l *Instance) Log(level string, message interface{}) {
	if lvl, err := logrus.ParseLevel(level); err == nil {
		if message != nil {
			l.logger.Log(lvl, message)
		}
	} else {
		l.logger.Fatal(err.Error())
	}
}

//Default logger you don't have to use this logger if you want to use other logger libraries implement logger interface
func DefaultLogger() ILogger {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02T15:04:05.999999Z",
	})
	return &Instance{logger: logrus.StandardLogger()}
}

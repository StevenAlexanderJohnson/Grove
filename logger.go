package grove

import (
	"log"
	"os"
)

type ILogger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type DefaultLogger struct {
	logger *log.Logger
}

func (l *DefaultLogger) Log(v ...interface{}) {
	l.logger.Println(v...)
}
func (l *DefaultLogger) Logf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}
func (l *DefaultLogger) Info(v ...interface{}) {
	v = append([]interface{}{"INFO:"}, v...)
	l.logger.Println(v...)
}
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.logger.Printf("INFO: "+format, v...)
}
func (l *DefaultLogger) Error(v ...interface{}) {
	v = append([]interface{}{"ERROR:"}, v...)
	l.logger.Println(v...)
}
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("ERROR: "+format, v...)
}
func (l *DefaultLogger) Debug(v ...interface{}) {
	v = append([]interface{}{"DEBUG:"}, v...)
	l.logger.Println(v...)
}
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.logger.Printf("DEBUG: "+format, v...)
}
func (l *DefaultLogger) Warning(v ...interface{}) {
	v = append([]interface{}{"WARNING:"}, v...)
	l.logger.Println(v...)
}
func (l *DefaultLogger) Warningf(format string, v ...interface{}) {
	l.logger.Printf("WARNING: "+format, v...)
}
func (l *DefaultLogger) Trace(v ...interface{}) {
	v = append([]interface{}{"TRACE:"}, v...)
	l.logger.Println(v...)
}
func (l *DefaultLogger) Tracef(format string, v ...interface{}) {
	l.logger.Printf("TRACE: "+format, v...)
}
func (l *DefaultLogger) Fatal(v ...interface{}) {
	v = append([]interface{}{"FATAL:"}, v...)
	l.logger.Println(v...)
	os.Exit(1)
}
func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Printf("FATAL: "+format, v...)
	os.Exit(1)
}

func NewDefaultLogger() ILogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "grove: ", log.LstdFlags),
	}
}

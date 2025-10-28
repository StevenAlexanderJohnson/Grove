package grove

import (
	"log"
	"os"
)

// Interface that Grove uses for logging. It allows end users to define their
// own logger or bring their own by providing the required methods.
type ILogger interface {
	Log(v ...any)
	Logf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Warning(v ...any)
	Warningf(format string, v ...any)
	Trace(v ...any)
	Tracef(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}

// The default logger Grove will use for logging.
// It simply adds the logging level before the message.
type DefaultLogger struct {
	logger *log.Logger
}

// Logs the information with no level specified.
func (l *DefaultLogger) Log(v ...any) {
	l.logger.Println(v...)
}

// Logs the message allowing you to format the string.
func (l *DefaultLogger) Logf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

// Logs the information prepended with "INFO".
func (l *DefaultLogger) Info(v ...any) {
	v = append([]any{"INFO:"}, v...)
	l.logger.Println(v...)
}

// Logs the message prepended with 'INFO' that allows you to format the string.
func (l *DefaultLogger) Infof(format string, v ...any) {
	l.logger.Printf("INFO: "+format, v...)
}

// Logs the message prepended with 'ERROR'.
func (l *DefaultLogger) Error(v ...any) {
	v = append([]any{"ERROR:"}, v...)
	l.logger.Println(v...)
}

// Logs the message prepended with 'ERROR' that allows you to format the string.
func (l *DefaultLogger) Errorf(format string, v ...any) {
	l.logger.Printf("ERROR: "+format, v...)
}

// Logs the message prepended with 'DEBUG'.
func (l *DefaultLogger) Debug(v ...any) {
	v = append([]any{"DEBUG:"}, v...)
	l.logger.Println(v...)
}

// Logs the message prepended with 'DEBUG' that allows you to format the string.
func (l *DefaultLogger) Debugf(format string, v ...any) {
	l.logger.Printf("DEBUG: "+format, v...)
}

// Logs the message prepended with 'WARNING'.
func (l *DefaultLogger) Warning(v ...any) {
	v = append([]any{"WARNING:"}, v...)
	l.logger.Println(v...)
}

// Logs the message prepended with 'WARNING' that allows you to format the string.
func (l *DefaultLogger) Warningf(format string, v ...any) {
	l.logger.Printf("WARNING: "+format, v...)
}

// Logs the message prepended with 'TRACE'.
func (l *DefaultLogger) Trace(v ...any) {
	v = append([]any{"TRACE:"}, v...)
	l.logger.Println(v...)
}

// Logs the message prepended with 'TRACE' that allows you to format the string.
func (l *DefaultLogger) Tracef(format string, v ...any) {
	l.logger.Printf("TRACE: "+format, v...)
}

// Logs the message prepended with 'FATAL'. It will also exit the application with code 1.
func (l *DefaultLogger) Fatal(v ...any) {
	v = append([]any{"FATAL:"}, v...)
	l.logger.Println(v...)
	os.Exit(1)
}

// Logs the message prepended with 'FATAL' that allows you to format the string.
// It will also exit the application with code 1.
func (l *DefaultLogger) Fatalf(format string, v ...any) {
	l.logger.Printf("FATAL: "+format, v...)
	os.Exit(1)
}

// Initializes the Default logger and prepends the `appName` to all log methods.
func NewDefaultLogger(appName string) ILogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, appName+": ", log.LstdFlags),
	}
}

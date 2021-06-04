package logger

import "fmt"

// Info logs INFO messages. stdout flag indicates if message is to be written to stdout in addition to log.
func Info(stdout bool, msg string) {
	logInfo(loggersMap.getLogger(), stdout, msg)
}

// Infof logs INFO messages. stdout flag indicates if message is to be written to stdout in addition to log.
func Infof(stdout bool, msg string, args ...interface{}) {
	Info(stdout, fmt.Sprintf(msg, args...))
}

// Debug logs DEBUG messages. stdout flag indicates if message is to be written to stdout in addition to log.
func Debug(stdout bool, msg string) {
	logDebug(loggersMap.getLogger(), stdout, msg)
}

// Debugf logs DEBUG messages. stdout flag indicates if message is to be written to stdout in addition to log.
func Debugf(stdout bool, msg string, args ...interface{}) {
	Debug(stdout, fmt.Sprintf(msg, args...))
}

//nolint:gochecknoglobals
// Package logger provides logging functionality.
// It only provides debug and info logging levels, following the recommendation at:
// https://dave.cheney.net/2015/11/05/lets-talk-about-logging
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/agilepathway/gauge-confluence/util"
	"github.com/natefinch/lumberjack"
	logging "github.com/op/go-logging"
)

const (
	gaugeConfluenceModuleID    = "GaugeConfluence"
	logsDirectory              = "logs_directory"
	logs                       = "logs"
	gaugeConfluenceLogFileName = "gaugeconfluence.log"
)

var level logging.Level
var initialized bool
var loggersMap logCache
var fileLogFormat = logging.MustStringFormatter("%{time:02-01-2006 15:04:05.000} [%{module}] [%{level}] %{message}")
var fileLoggerLeveled logging.LeveledBackend

// ActiveLogFile log file represents the file which will be used for the backend logging
var ActiveLogFile string

type logCache struct {
	mutex   sync.RWMutex
	loggers map[string]*logging.Logger
}

// getLogger gets logger for the gaugeConfluencePlugin.
func (l *logCache) getLogger() *logging.Logger {
	if !initialized {
		return nil
	}

	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if _, ok := l.loggers[gaugeConfluenceModuleID]; !ok {
		l.mutex.RUnlock()
		l.addLogger(gaugeConfluenceModuleID)
		l.mutex.RLock()
	}

	return l.loggers[gaugeConfluenceModuleID]
}

func (l *logCache) addLogger(module string) {
	logger := logging.MustGetLogger(module)
	logger.SetBackend(fileLoggerLeveled)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.loggers[module] = logger
}

// Initialize logger with given level
func Initialize(logLevel string) {
	level = loggingLevel(logLevel)
	ActiveLogFile = getLogFile()

	initFileLoggerBackend()

	loggersMap = logCache{loggers: make(map[string]*logging.Logger)}
	loggersMap.addLogger(gaugeConfluenceModuleID)

	initialized = true
}

func logInfo(logger *logging.Logger, stdout bool, msg string) {
	if level >= logging.INFO {
		write(stdout, msg, os.Stdout)
	}

	if !initialized {
		return
	}

	logger.Infof(msg)
}

func logDebug(logger *logging.Logger, stdout bool, msg string) {
	if level >= logging.DEBUG {
		write(stdout, msg, os.Stdout)
	}

	if !initialized {
		return
	}

	logger.Debugf(msg)
}

func write(stdout bool, msg string, writer io.Writer) {
	if stdout {
		fmt.Fprintln(writer, msg)
	}
}

func initFileLoggerBackend() {
	var backend = createFileLogger(ActiveLogFile, 10)
	fileFormatter := logging.NewBackendFormatter(backend, fileLogFormat)
	fileLoggerLeveled = logging.AddModuleLevel(fileFormatter)
	fileLoggerLeveled.SetLevel(logging.DEBUG, "")
}

//nolint:gomnd
var createFileLogger = func(name string, size int) logging.Backend {
	return logging.NewLogBackend(&lumberjack.Logger{
		Filename:   name,
		MaxSize:    size, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}, "", 0)
}

func addLogsDirPath(logFileName string) string {
	customLogsDir := os.Getenv(logsDirectory)
	if customLogsDir == "" {
		return filepath.Join(logs, logFileName)
	}

	return filepath.Join(customLogsDir, logFileName)
}

func getLogFile() string {
	logDirPath := addLogsDirPath(gaugeConfluenceLogFileName)
	if filepath.IsAbs(logDirPath) {
		return logDirPath
	}

	return filepath.Join(getProjectRoot(), logDirPath)
}

func getProjectRoot() string {
	return util.GetProjectRoot()
}

func loggingLevel(logLevel string) logging.Level {
	if logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug":
			return logging.DEBUG
		case "info":
			return logging.INFO
		case "warning":
			return logging.WARNING
		case "error":
			return logging.ERROR
		}
	}

	return logging.INFO
}

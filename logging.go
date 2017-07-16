package log

import (
	"github.com/op/go-logging"
	"io"
	"os"
	"strings"
)

// default log level for file logger
const defaultLogLevel = "info"

// Logger should be used by all modules to write logs
var Logger *logging.Logger

// coloredFormat would be used for pretty console output
var coloredFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

// uncoloredFormat would be used for log file output
var uncoloredFormat = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{shortfunc} %{level:.4s} %{message}`,
)

var loggingLevels = map[string]logging.Level{
	"debug":    logging.DEBUG,
	"info":     logging.INFO,
	"notice":   logging.NOTICE,
	"warning":  logging.WARNING,
	"error":    logging.ERROR,
	"critical": logging.CRITICAL,
}

func logLevelFromString(levelName string) logging.Level {
	// this function maps provided string logging level to actual logging
	// level type
	if level, found := loggingLevels[strings.ToLower(levelName)]; found {
		return level
	}
	return loggingLevels[defaultLogLevel]
}

// InitializeLogger sets up Logger
func InitializeLogger(logRotate io.Writer, consoleDebug bool, logLevel string) {
	Logger = logging.MustGetLogger("")

	// initialize console backend
	consoleBackend := logging.NewLogBackend(os.Stdout, "", 0)
	consoleFormatter := logging.NewBackendFormatter(
		consoleBackend, coloredFormat)
	consoleLeveledFormatter := logging.AddModuleLevel(consoleFormatter)
	var consoleLoggingLevel logging.Level
	if consoleDebug {
		consoleLoggingLevel = logging.DEBUG
	} else {
		consoleLoggingLevel = logging.INFO
	}
	consoleLeveledFormatter.SetLevel(consoleLoggingLevel, "")

	fileBackend := logging.NewLogBackend(logRotate, "", 0)
	fileFormatter := logging.NewBackendFormatter(
		fileBackend, uncoloredFormat)
	fileLeveledFormatter := logging.AddModuleLevel(fileFormatter)
	fileLeveledFormatter.SetLevel(logLevelFromString(logLevel), "")

	// set 2 backends
	logging.SetBackend(consoleLeveledFormatter, fileLeveledFormatter)
}

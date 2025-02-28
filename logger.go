package emu

import (
	"io"
	"log"
	"os"
)

var (
	DebugLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	logFile       io.Writer = os.Stdout
)

func enableLogger(l *log.Logger) {
	l.SetOutput(logFile)
}
func disableLogger(l *log.Logger) {
	l.SetOutput(io.Discard)
}

func initLog(file io.Writer, l LogLevel) {
	logFile = file
	var debugFile, infoFile, warningFile, errorFile io.Writer = io.Discard, io.Discard, io.Discard, io.Discard
	switch l {
	case LOG_ALL:
		debugFile, infoFile, warningFile, errorFile = file, file, file, file
	case LOG_ERROR:
		errorFile = file
	case LOG_WARNING:
		warningFile, errorFile = file, file
	case LOG_INFO:
		infoFile, warningFile, errorFile = file, file, file
	case LOG_OFF:
	default:
		debugFile, infoFile, warningFile, errorFile = io.Discard, io.Discard, io.Discard, io.Discard
	}
	DebugLogger = log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(warningFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

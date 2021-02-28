package plog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	stdLogger = New(allLog, "")

	logLevels = map[logLevel]string{
		LevelDebug: "Debug",
		LevelInfo:  "Info",
		LevelWarn:  "Warn",
		LevelError: "Error",
		LevelFatal: "Fatal",
	}
)

type (
	logLevel uint8

	logger struct {
		sysLogLevel logLevel
		fileOut     string
	}
)

const (
	allLog logLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	noLog
)

func New(sysLevel logLevel, filePath string) *logger {
	return &logger{sysLogLevel: sysLevel, fileOut: filePath}
}

func Logf(l logLevel, format string, a ...interface{}) {
	if l < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(l, format, a...)
}

func (lg *logger) writeLog(l logLevel, format string, a ...interface{}) {
	const timeFormat = "2006-01-02 15:04:05.000"
	pc, file, line, _ := runtime.Caller(2)

	fileNames := strings.Split(file, "/")
	fileName := fileNames[len(fileNames)-1]
	fc := runtime.FuncForPC(pc).Name()
	fcStrs := strings.Split(fc, "/")
	fcStr := fcStrs[len(fcStrs)-1]

	buf := fmt.Sprintf("[%s][%s][%s,%d][%s] ", time.Now().Format(timeFormat), logLevels[l], fileName, line, fcStr)
	buf += fmt.Sprintf(format+"\n", a...)
	if lg.fileOut != "" {
		fd, err := os.OpenFile(lg.fileOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return
		}
		fd.Write([]byte(buf))
		fd.Close()
	} else {
		os.Stdout.Write([]byte(buf))
	}
}

func Fatalf(format string, a ...interface{}) {
	if LevelFatal < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(LevelFatal, format, a...)
}

func Errorf(format string, a ...interface{}) {
	if LevelError < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(LevelError, format, a...)
}

func Warnf(format string, a ...interface{}) {
	if LevelWarn < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(LevelWarn, format, a...)
}

func Infof(format string, a ...interface{}) {
	if LevelInfo < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(LevelInfo, format, a...)
}

func Debugf(format string, a ...interface{}) {
	if LevelDebug < stdLogger.sysLogLevel {
		return
	}
	stdLogger.writeLog(LevelDebug, format, a...)
}

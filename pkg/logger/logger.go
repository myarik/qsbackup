package logger

import (
	"io"
	"log"
	"bytes"
	"strings"
)

const (
	DEBUG    = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

var logLevelName = [...]string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

type Log struct {
	logLevel int
	logger   *log.Logger
}

func New(out io.Writer, level int) *Log {
	return &Log{
		logLevel: level,
		logger:   log.New(out, "", log.LstdFlags),
	}
}

func (l Log) formatMessage(level, message string) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.ToUpper(level))
	buffer.WriteString("]: ")
	buffer.WriteString(message)
	return buffer.String()
}
func (l Log) log(logLevel int, message string) {
	if logLevel == CRITICAL {
		l.logger.Panicln(l.formatMessage(logLevelName[logLevel], message))
	} else if logLevel >= l.logLevel {
		l.logger.Println(l.formatMessage(logLevelName[logLevel], message))
	}
}

func (l Log) Debug(message string) {
	l.log(DEBUG, message)
}

func (l Log) Info(message string) {
	l.log(INFO, message)
}

func (l Log) Warning(message string) {
	l.log(WARNING, message)
}
func (l Log) Error(message string) {
	l.log(ERROR, message)
}
func (l Log) Critical(message string) {
	l.log(CRITICAL, message)
}

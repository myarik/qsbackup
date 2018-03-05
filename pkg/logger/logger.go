package logger

import (
	"io"
	"log"
	"bytes"
)

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

func (l Log) log(logLevel, message string) {
	var buffer bytes.Buffer
	buffer.WriteString(logLevel)
	buffer.WriteString(": ")
	buffer.WriteString(message)
	l.logger.Println(buffer.String())
}

func (l Log) Debug(message string) {
	l.log("DEBUG", message)
}

func (l Log) Info(message string) {
	l.log("INFO", message)
}

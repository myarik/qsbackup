/*
Usage

	import log "github.com/myarik/qsbackup/pkg/logger"

	func main()  {
		logger := log.New(os.Stdout, 0)
		logger.Debug("Test") 		// 2018/03/14 18:11:34 [DEBUG]: Test
		logger.Info("Test") 		// 2018/03/14 18:11:34 [WARNING]: Test
		logger.Warning("Test")      // 2018/03/14 18:11:34 [WARNING]: Test
		logger.Error("Test")        // 2018/03/14 18:11:34 [ERROR]: Test
		logger.Critical("Panic")    // 2018/03/14 18:11:34 [CRITICAL]: Panic  A program crashes
	}
 */
package logger

import (
	"io"
	"log"
	"bytes"
	"strings"
	"os"
	"io/ioutil"
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

type CloseStdout interface {
	Close() error
}

type Log struct {
	logLevel        int
	logger          *log.Logger
	closableOutputs []CloseStdout
}

func New(out io.Writer, level int) *Log {
	return &Log{
		logLevel: level,
		logger:   log.New(out, "", log.LstdFlags),
	}
}

// Setup a logger
func Init(logFile string, debugMode bool) (*Log, error) {
	if logFile != "" {
		fileStdout, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		var logger *Log
		if debugMode == true {
			logger = New(io.MultiWriter(fileStdout, os.Stdout), 0)
		} else {
			logger = New(fileStdout, 1)
		}
		logger.closableOutputs = append(logger.closableOutputs, fileStdout)
		return logger, nil
	} else if debugMode == true {
		return New(os.Stdout, 0), nil
	}
	return New(ioutil.Discard, 1), nil
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

func (l Log) Close() error {
	for _, item := range l.closableOutputs {
		item.Close()
	}
	return nil
}

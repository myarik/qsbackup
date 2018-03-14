package logger

import (
	"testing"
	"bytes"
	"time"
	"fmt"
)

const pattern = "%d/%.2d/%.2d %.2d:%.2d:%.2d %s: %s\n"

func TestOutput(t *testing.T) {
	const testMgs = "Lorem Ipsum is simply dummy text of the printing and typesetting industry."
	var outBuffer bytes.Buffer
	logger := New(&outBuffer, 0)
	now := time.Now()
	expect := fmt.Sprintf(
		pattern,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), "[INFO]", testMgs)
	logger.Info(testMgs)
	if outBuffer.String() != expect {
		t.Errorf("log output should match %q is %q", expect, outBuffer.String())
	}
}

func TestDebugOutput(t *testing.T) {
	const testMgs = "Lorem Ipsum is simply dummy text of the printing and typesetting industry."
	var outBuffer bytes.Buffer
	logger := New(&outBuffer, 1)
	logger.Debug(testMgs)
	if outBuffer.String() != "" {
		t.Errorf("log output should be empty", )
	}
	now := time.Now()
	expect := fmt.Sprintf(
		pattern,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), "[ERROR]", testMgs)
	logger.Error(testMgs)
	if outBuffer.String() != expect {
		t.Errorf("log output should match %q is %q", expect, outBuffer.String())
	}
}

func TestPanic(t *testing.T) {
	const testMgs = "Lorem Ipsum is simply dummy text of the printing and typesetting industry."
	var outBuffer bytes.Buffer
	logger := New(&outBuffer, 0)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	// The following is the code under test
	logger.Critical(testMgs)
}

package log

import (
	"testing"
)

func TestLoggerDebug(t *testing.T) {
	logger := New("robokiller-ivr", "1.0")
	logger.Set("foo", "bar")
	logger.Debug("debug message")
}

func TestLoggerInfo(t *testing.T) {
	logger := New("robokiller-ivr", "1.0")
	logger.Set("foo", "bar")
	logger.Info("info message")
}

func TestLoggerError(t *testing.T) {
	logger := New("robokiller-ivr", "1.0")
	logger.Set("foo", "bar")
	logger.Error("error message")
}

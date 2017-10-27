package logger

import (
	"time"
	"fmt"
	"encoding/json"
	"runtime"
	"bytes"
	"os"
	"io"
)

type severity int

const (
	DEBUG severity = iota
	INFO
	WARN
	ERROR
)

var LogLevelName = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

var LogLevelValue = map[string]severity{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

type Fields map[string]string

type ServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

type ReportLocation struct {
	FilePath     string `json:"filePath"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
}

type Context struct {
	Data           Fields          `json:"data,omitempty"`
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

type Payload struct {
	Severity       string          `json:"severity"`
	EventTime      string          `json:"eventTime"`
	Caller         string          `json:"caller,omitempty"`
	Message        string          `json:"message"`
	ServiceContext *ServiceContext `json:"serviceContext"`
	Context        *Context        `json:"context,omitempty"`
	Stacktrace     string          `json:"stacktrace,omitempty"`
}

type Log struct {
	Payload *Payload
	writer io.Writer
}

func New() *Log {
	if os.Getenv("SERVICE") == "" || os.Getenv("VERSION") == "" {
		fmt.Errorf("cannot instantiate the logger, make sure the SERVICE and VERSION environment vars are set correctly")
	}

	return &Log{
		Payload: &Payload{
			ServiceContext: &ServiceContext{
				Service: os.Getenv("SERVICE"),
				Version: os.Getenv("VERSION"),
			},
		},
		writer: os.Stdout,
	}
}

func (l *Log) SetWriter(w io.Writer) *Log {
	l.writer = w
	return l
}

func (l *Log) set(key, val string) {
	if l.Payload.Context == nil {
		l.Payload.Context = &Context{
			Data: Fields{},
		}
	}

	l.Payload.Context.Data[key] = val
}

func (l *Log) log(severity, message string) {
	l.Payload = &Payload{
		Severity: severity,
		EventTime: time.Now().Format(time.RFC3339),
		Message: message,
		ServiceContext: l.Payload.ServiceContext,
		Context: l.Payload.Context,
		Stacktrace: l.Payload.Stacktrace,
	}

	payload, ok := json.Marshal(l.Payload)
	if ok != nil {
		fmt.Errorf("cannot marshal payload: %s", ok.Error())
	}

	fmt.Fprintln(l.writer, string(payload))
}

// Checks whether the specified log level is valid in the current environment
func isValidLogLevel(logLevel severity) bool {
	curLogLev, ok := LogLevelValue[os.Getenv("LOG_LEVEL")]
	if !ok {
		fmt.Errorf("the LOG_LEVEL environment variable is not set or has an incorrect value")
	}

	return curLogLev <= logLevel
}

func (l *Log) With(fields Fields) *Log {
	if os.Getenv("SERVICE") == "" || os.Getenv("VERSION") == "" {
		fmt.Errorf("cannot instantiate the logger, make sure the SERVICE and VERSION environment vars are set correctly")
	}

	log := &Log{
		Payload: &Payload{
			ServiceContext: &ServiceContext{
				Service: os.Getenv("SERVICE"),
				Version: os.Getenv("VERSION"),
			},
		},
		writer: os.Stdout,
	}

	for k, v := range fields {
		log.set(k, v)
	}

	return log
}

func (l *Log) Debug(message string) {
	if !isValidLogLevel(DEBUG) {
		return
	}

	l.log(LogLevelName[DEBUG], message)
}

func (l *Log) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

func (l *Log) Metric(message string) {
	if !isValidLogLevel(INFO) {
		return
	}

	l.log(LogLevelName[INFO], message)
}

func (l *Log) Info(message string) {
	if !isValidLogLevel(INFO) {
		return
	}

	l.log(LogLevelName[INFO], message)
}

func (l *Log) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

func (l *Log) Warn(message string) {
	if !isValidLogLevel(WARN) {
		return
	}

	l.log(LogLevelName[WARN], message)
}

func (l *Log) Warnf(message string, args ...interface{}) {
	l.Warn(fmt.Sprintf(message, args...))
}

func (l *Log) Error(message string) {
	buffer := make([]byte, 1024)
	runtime.Stack(buffer, false)
	_, file, line, _ := runtime.Caller(1)

	l.Payload = &Payload{
		ServiceContext: l.Payload.ServiceContext,
		Context: &Context{
			Data: l.Payload.Context.Data,
			ReportLocation: &ReportLocation{
				FilePath: file,
				FunctionName: "unknown",
				LineNumber: line,
			},
		},
		Stacktrace: string(bytes.Trim(buffer, "\x00")),
	}

	l.log(LogLevelName[ERROR], message)
}

func (l *Log) Errorf(message string, args ...interface{}) {
	l.Error(fmt.Sprintf(message, args...))
}

package log

import (
	"time"
	"fmt"
	"encoding/json"
	"runtime"
	"bytes"
)

const (
	SEVERITY_DEBUG = "DEBUG"
	SEVERITY_INFO  = "INFO"
	SEVERITY_WARN  = "WARN"
	SEVERITY_ERROR = "ERROR"
)

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
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

type Payload struct {
	Severity       string             `json:"severity"`
	EventTime      string             `json:"eventTime"`
	Caller         string             `json:"caller,omitempty"`
	Message        string             `json:"message"`
	Data           map[string]string  `json:"data,omitempty"`
	ServiceContext *ServiceContext    `json:"serviceContext"`
	Context        *Context           `json:"context,omitempty"`
	Stacktrace     string             `json:"stacktrace,omitempty"`
}

type Logger struct {
	Payload *Payload
}

func New(service, version string) *Logger {
	return &Logger{
		Payload: &Payload{
			ServiceContext: &ServiceContext{
				Service: service,
				Version: version,
			},
		},
	}
}

func (l *Logger) Set(key, val string) {
	if l.Payload.Data == nil {
		l.Payload.Data = map[string]string{}
	}

	l.Payload.Data[key] = val
}

func (l *Logger) log(severity, message string) {
	l.Payload = &Payload{
		Severity: severity,
		EventTime: time.Now().Format(time.RFC3339),
		Message: message,
		Data: l.Payload.Data,
		ServiceContext: l.Payload.ServiceContext,
		Context: l.Payload.Context,
		Stacktrace: l.Payload.Stacktrace,
	}

	line, ok := json.Marshal(l.Payload)
	if ok != nil {
		fmt.Errorf("cannot marshal payload: %s", ok.Error())
	}

	fmt.Println(string(line))

	// Unset the current payload data
	l.Payload.Data = nil
}

func (l *Logger) Debug(message string) {
	l.log(SEVERITY_DEBUG, message)
}

func (l *Logger) Info(message string) {
	l.log(SEVERITY_INFO, message)
}

func (l *Logger) Warn(message string) {
	l.log(SEVERITY_WARN, message)
}

func (l *Logger) Error(message string) {
	buffer := make([]byte, 1024)
	runtime.Stack(buffer, false)
	_, file, line, _ := runtime.Caller(1)
	l.Payload = &Payload{
		ServiceContext: l.Payload.ServiceContext,
		Context: &Context{
			ReportLocation: &ReportLocation{
				FilePath: file,
				FunctionName: "unknown",
				LineNumber: line,
			},
		},
		Stacktrace: string(bytes.Trim(buffer, "\x00")),
	}

	l.log(SEVERITY_ERROR, message)
}
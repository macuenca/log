package logger

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
)

const OUTFILE = "out.json"

func setEnv() {
	os.Setenv("SERVICE", "robokiller-ivr")
	os.Setenv("VERSION", "1.0")
}

func createOutFile() *os.File {
	// Delete file first if exists
	os.Remove(OUTFILE)

	file, err := os.Create(OUTFILE)
	if err != nil {
		panic("Unable to create test file")
	}

	return file
}

func compareWithOutFile(expected string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	return strings.TrimRight(string(data), "\n") == expected
}

func outFileContains(substring string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	fileData := strings.TrimRight(string(data), "\n")
	return strings.Contains(fileData, substring)
}

func TestLoggerDebugWithImplicitContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key": "value",
		"function" : "TestLoggerDebug",
	}).SetWriter(file)

	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}


func TestLoggerDebugWithExplicitContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key": "value",
		"function" : "TestLoggerDebug",
	}).SetWriter(file)

	log.With(Fields{"function": "TestLoggerDebug"}).SetWriter(file).Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerDebugWithoutContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().SetWriter(file)

	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerDebugfWithoutContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().SetWriter(file)

	param := "with param"
	log.Debugf("debug message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerMetric(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().SetWriter(file)

	log.Metric("custom_metric")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"custom_metric\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerInfo(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerInfo",
	}).SetWriter(file)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerInfof(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerInfo",
	}).SetWriter(file)

	param := "with param"
	log.Infof("info message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerError(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerError",
	}).SetWriter(file)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"}") {
		t.Errorf("output file %s does not contain the context", OUTFILE)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", OUTFILE)
	}
}

func TestLoggerErrorf(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerError",
	}).SetWriter(file)

	param := "with param"
	log.Errorf("error message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}
}

func TestLoggerInfoWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"function": "TestLoggerInfo",
		"key":"value",
		"package": "logger",
	}).SetWriter(file)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\",\"package\":\"logger\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerErrorWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New().With(Fields{
		"function": "TestLoggerError",
		"key":"value",
		"package": "logger",
	}).SetWriter(file)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\",\"package\":\"logger\"}") {
		t.Errorf("output file %s does not contain the context", OUTFILE)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", OUTFILE)
	}
}

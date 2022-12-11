// Package logging provides an abstraction about the logging setup.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Logger provides an interface for differnet types of loggers.
type Logger struct {
	Name       string
	Writer     *os.File
	Timeformat string
	Level      string
	Formatter  Formatter
}

// Formatter is the interface any log formatter needs to implement.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// JSONFormatter formats log message in JSON so they can be easily parsed.
type JSONFormatter struct{}

// TextFormatter formats logs in plain text.
type TextFormatter struct{}

// Format for JSONFormatter formats the log item as a dictionary.
func (j *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	logItem := make(map[string]interface{}, 0)
	logItem["app"] = entry.Logger.Name
	logItem["time"] = time.Now().Format(entry.Logger.Timeformat)

	for i := 0; i < len(entry.Data); i += 2 {
		logItem[fmt.Sprintf("%v", entry.Data[i])] = entry.Data[i+1]
	}

	logItemJSON, err := json.Marshal(logItem)
	if err != nil {
		return nil, err
	}

	return logItemJSON, nil
}

// Format for TextFormatter formats the log item as plain text.
func (j *TextFormatter) Format(entry *Entry) ([]byte, error) {
	logItem := fmt.Sprintf("time: %s app: %s", time.Now().Format(entry.Logger.Timeformat), entry.Logger.Name)

	for i := 0; i < len(entry.Data); i += 2 {
		logItem = fmt.Sprintf("%s %s: %v", logItem, fmt.Sprintf("%v", entry.Data[i]), entry.Data[i+1])
	}

	return []byte(logItem), nil
}

// Entry provides a log entry.
type Entry struct {
	Logger Logger
	Data   []interface{}
	Level  string
}

const (
	defaultTimeFormat = "2006-01-02 15:04:05" // defaultTimeFormat represents the default time format.
	defaultLogLevel   = "info"                // defaultLogLevel represents the default log level used.
)

// Info provides an information level logging capability.
func (e *Entry) Info(kvs ...interface{}) {
	if len(kvs)%2 != 0 {
		e.Data = nil
		e.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	if e.Logger.Level != "info" {
		return
	}

	e.Data = kvs
	j, err := e.Logger.Formatter.Format(e)

	if err != nil {
		e.Data = nil
		e.Error(fmt.Errorf("unable to format message"), "unable to log")

		return
	}

	e.Logger.Write(e, j)
}

// Error provides an error level logging capability.
func (e *Entry) Error(err error, msg string, kvs ...interface{}) {
	if e.Logger.Level != "error" {
		return
	}

	if len(kvs)%2 != 0 {
		e.Data = nil
		e.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	kvs = append(kvs, err.Error(), msg)
	e.Data = kvs
	j, err := e.Logger.Formatter.Format(e)

	if err != nil {
		e.Data = nil
		e.Error(fmt.Errorf("unable to format message"), "unable to log")

		return
	}

	e.Logger.Write(e, j)
}

// NewEntry provides a logging entry.
func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: *logger,
	}
}

// Write write out the log line.
func (l Logger) Write(entry *Entry, data []byte) {
	fmt.Fprintf(entry.Logger.Writer, "%s\n", string(data))
	entry.Data = nil
}

// New returns a new logging to defaulting to stadard out.
func New(name string) *Logger {
	return &Logger{
		Name:       name,
		Writer:     os.Stdout,
		Timeformat: defaultTimeFormat,
		Level:      defaultLogLevel,
		Formatter:  &TextFormatter{},
	}
}

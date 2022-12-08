// Package logging provides an abstraction about the logging setup.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

// Logger provides an interface for differnet types of loggers.
type Logger interface {
	Write(entry *Entry)
}

// JSONLogger handles logging all data.
type JSONLogger struct {
	Name       string
	Writer     *os.File
	Timeformat string
}

// TabLogger handles logging all data.
type TabLogger struct {
	Name       string
	Writer     *os.File
	Timeformat string
}

// PlainLogger handles logging all data.
type PlainLogger struct {
	Name       string
	Writer     *os.File
	Timeformat string
}

// Entry provides a log entry.
type Entry struct {
	Logger Logger
	Data   []interface{}
}

const (
	defaultTimeFormat = "2006-01-02 15:04:05" // defaultTimeFormat represents the default time format.
)

// Info provides an information level logging capability.
func (e *Entry) Info(kvs ...interface{}) {
	if len(e.Data)%2 != 0 {
		e.Data = nil
		e.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	e.Data = kvs
	e.Logger.Write(e)
}

// Error provides an error level logging capability.
func (e *Entry) Error(err error, msg string, kvs ...interface{}) {
	e.Data = append(e.Data, err.Error(), msg)
	e.Info(e.Data...)
}

// NewEntry provides a logging entry.
func NewEntry(logger Logger) *Entry {
	return &Entry{
		Logger: logger,
	}
}

// Write write out the log line.
func (l JSONLogger) Write(entry *Entry) {
	logItem := make(map[string]interface{}, 0)
	logItem["app"] = entry.Logger.(*JSONLogger).Name
	logItem["time"] = time.Now().Format(entry.Logger.(*JSONLogger).Timeformat)

	if len(entry.Data)%2 != 0 {
		entry.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	for i := 0; i < len(entry.Data); i += 2 {
		logItem[fmt.Sprintf("%v", entry.Data[i])] = entry.Data[i+1]
	}

	logItemJSON, err := json.Marshal(logItem)
	if err != nil {
		entry.Error(err, "unable to log")

		return
	}

	fmt.Fprintf(entry.Logger.(*JSONLogger).Writer, "%s\n", string(logItemJSON))
	entry.Data = nil
}

// NewJSONLogger returns a new logging to defaulting to stadard out.
func NewJSONLogger(name string, timeformat string) *JSONLogger {
	if timeformat == "" {
		timeformat = defaultTimeFormat
	}

	return &JSONLogger{
		Name:       name,
		Writer:     os.Stdout,
		Timeformat: timeformat,
	}
}

// Write write out the log line.
func (l TabLogger) Write(entry *Entry) {
	w := tabwriter.NewWriter(l.Writer, 1, 1, 1, ' ', 0)
	logItem := fmt.Sprintf("time: %s\tapp: %s", time.Now().Format(entry.Logger.(*TabLogger).Timeformat), entry.Logger.(*TabLogger).Name)

	if len(entry.Data)%2 != 0 {
		entry.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	for i := 0; i < len(entry.Data); i += 2 {
		logItem = fmt.Sprintf("%s\t%s:\t%v", logItem, fmt.Sprintf("%v", entry.Data[i]), entry.Data[i+1])
	}

	fmt.Fprintf(w, "%s\n", logItem)

	err := w.Flush()

	if err != nil {
		entry.Error(err, "unable to log")

		return
	}

	entry.Data = nil
}

// NewTabLogger returns a new logging to defaulting to stadard out.
func NewTabLogger(name string, timeformat string) *TabLogger {
	if timeformat == "" {
		timeformat = defaultTimeFormat
	}

	return &TabLogger{
		Name:       name,
		Writer:     os.Stdout,
		Timeformat: timeformat,
	}
}

// Write write out the log line.
func (l PlainLogger) Write(entry *Entry) {
	logItem := fmt.Sprintf("time: %s app: %s", time.Now().Format(entry.Logger.(*PlainLogger).Timeformat), entry.Logger.(*PlainLogger).Name)

	if len(entry.Data)%2 != 0 {
		entry.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	for i := 0; i < len(entry.Data); i += 2 {
		logItem = fmt.Sprintf("%s %s: %v", logItem, fmt.Sprintf("%v", entry.Data[i]), entry.Data[i+1])
	}

	fmt.Fprintf(entry.Logger.(*PlainLogger).Writer, "%s\n", logItem)

	entry.Data = nil
}

// NewPlainLogger returns a new logging to defaulting to stadard out.
func NewPlainLogger(name string, timeformat string) *PlainLogger {
	if timeformat == "" {
		timeformat = defaultTimeFormat
	}

	return &PlainLogger{
		Name:       name,
		Writer:     os.Stdout,
		Timeformat: timeformat,
	}
}

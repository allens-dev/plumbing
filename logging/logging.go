// Package logging provides an abstraction about the logging setup.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONLogger handles logging all data.
type JSONLogger struct {
	Name   string
	Writer *os.File
}

// Info provides an information level logging capability.
func (l JSONLogger) Info(kvs ...interface{}) {
	logItem := make(map[string]interface{}, 0)
	logItem["component"] = l.Name

	if len(kvs)%2 != 0 {
		l.Error(fmt.Errorf("need an even amount of kvs"), "unable to log")

		return
	}

	for i := 0; i < len(kvs); i += 2 {
		logItem[fmt.Sprintf("%v", kvs[i])] = kvs[i+1]
	}

	logItemJSON, err := json.Marshal(logItem)
	if err != nil {
		l.Error(err, "unable to log")

		return
	}

	fmt.Fprintf(l.Writer, "%s\n", string(logItemJSON))
}

// Error provides an error level logging capability.
func (l JSONLogger) Error(err error, msg string, kvs ...interface{}) {
	kvs = append(kvs, msg, err)
	l.Info(kvs...)
}

// NewJSONLogger returns a new logging to defaulting to stadard out.
func NewJSONLogger(name string) *JSONLogger {
	return &JSONLogger{
		Name:   name,
		Writer: os.Stdout,
	}
}

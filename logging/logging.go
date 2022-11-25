// Package logging provides an abstraction about the logging setup.
package logging

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONLogger struct {
	Name   string
	Writer *os.File
}

func (l JSONLogger) Info(kvs ...interface{}) {
	logItem := make(map[string]interface{}, 0)
	logItem["component"] = l.Name

	for i := 0; i < len(kvs); i += 2 {
		logItem[fmt.Sprintf("%v", kvs[i])] = kvs[i+1]
	}

	logItemJson, err := json.Marshal(logItem)
	if err != nil {
		l.Error(err, "unable to log")
	}
	fmt.Fprintf(l.Writer, "%s\n", string(logItemJson))
}

func (l JSONLogger) Error(err error, msg string, kvs ...interface{}) {
	kvs = append(kvs, msg, err)
	l.Info(kvs...)
}

func NewJSONLogger(name string) *JSONLogger {
	return &JSONLogger{
		Name:   name,
		Writer: os.Stdout,
	}
}

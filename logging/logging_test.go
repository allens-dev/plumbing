// Package logging provides an abstraction about the logging setup.
package logging_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/allens-dev/plumbing/logging"
)

func TestNewJSONLogger(t *testing.T) {
	type args struct {
		name string
	}

	loggerName := "test-logging"

	testWant := logging.NewJSONLogger(loggerName)

	testArgs := args{
		name: loggerName,
	}

	tests := []struct {
		name string
		args args
		want *logging.JSONLogger
	}{
		{name: "test-new-logger", args: testArgs, want: testWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logging.NewJSONLogger(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONLogger_Error(t *testing.T) {
	type fields struct {
		Name   string
		Writer *os.File
	}

	testFields := fields{
		Name:   "test-logger",
		Writer: os.Stdout,
	}

	type args struct {
		err error
		msg string
		kvs []interface{}
	}

	testArgs := args{
		err: fmt.Errorf("test error"),
		msg: "error during test",
		kvs: make([]interface{}, 0),
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test-error-log", fields: testFields, args: testArgs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logging.JSONLogger{
				Name:   tt.fields.Name,
				Writer: tt.fields.Writer,
			}
			l.Error(tt.args.err, tt.args.msg, tt.args.kvs...)
		})
	}
}

func TestJSONLogger_Info(t *testing.T) {
	type fields struct {
		Name   string
		Writer *os.File
	}

	testFields := fields{
		Name:   "test-logger",
		Writer: os.Stdout,
	}

	type args struct {
		kvs []interface{}
	}

	testKVS := make([]interface{}, 0)
	testKVS = append(testKVS, "starting", "up and reading for tests")

	testKVS2 := make([]interface{}, 0)
	testKVS2 = append(testKVS2, "starting")

	testArgs := args{
		kvs: testKVS,
	}

	testArgs2 := args{
		kvs: testKVS2,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test-info-logging", fields: testFields, args: testArgs},
		{name: "test-info-logging-2", fields: testFields, args: testArgs2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logging.JSONLogger{
				Name:   tt.fields.Name,
				Writer: tt.fields.Writer,
			}
			l.Info(tt.args.kvs...)
		})
	}
}

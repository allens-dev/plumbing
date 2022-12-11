// Package logging provides an abstraction about the logging setup.
package logging_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/allens-dev/plumbing/logging"
)

const testLoggerName = "test-logging"

func TestNew(t *testing.T) {
	type args struct {
		name string
	}

	loggerName := testLoggerName

	testWant := logging.New(loggerName)

	testArgs := args{
		name: loggerName,
	}

	tests := []struct {
		name string
		args args
		want *logging.Logger
	}{
		{name: "test-new-logger", args: testArgs, want: testWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logging.New(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntry(t *testing.T) {
	type args struct {
		logger *logging.Logger
	}

	loggerName := testLoggerName

	testLogger := logging.New(loggerName)

	testArgs := args{
		logger: testLogger,
	}

	testWant := logging.NewEntry(testLogger)

	tests := []struct {
		name string
		args args
		want *logging.Entry
	}{
		{name: "test-new-entry", args: testArgs, want: testWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logging.NewEntry(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Error(t *testing.T) {
	type fields struct {
		Logger *logging.Logger
		Data   []interface{}
	}

	loggerName := testLoggerName

	testLogger := logging.New(loggerName)
	testLogger.Level = "error"

	testFields := fields{
		Logger: testLogger,
		Data:   make([]interface{}, 0),
	}

	type args struct {
		err error
		msg string
		kvs []interface{}
	}

	testArgs := args{
		err: fmt.Errorf("this is an error"),
		msg: "this is why we error",
		kvs: make([]interface{}, 0),
	}

	testKVS := make([]interface{}, 0)
	testKVS = append(testKVS, "jello")

	testArgs2 := args{
		err: fmt.Errorf("this is an error"),
		msg: "this is why we error",
		kvs: testKVS,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test-error-write", fields: testFields, args: testArgs},
		{name: "test-error-write-odd-kvs", fields: testFields, args: testArgs2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &logging.Entry{
				Logger: *tt.fields.Logger,
				Data:   tt.fields.Data,
			}
			e.Error(tt.args.err, tt.args.msg, tt.args.kvs...)
		})
	}
}

func TestEntry_Info(t *testing.T) {
	type fields struct {
		Logger *logging.Logger
		Data   []interface{}
	}

	loggerName := testLoggerName

	testLogger := logging.New(loggerName)

	testLogger2 := logging.New(loggerName)
	testLogger2.Formatter = &logging.JSONFormatter{}

	testLogger3 := logging.New(loggerName)
	testLogger3.Formatter = &logging.JSONFormatter{}
	testLogger3.Level = "error"

	testFields := fields{
		Logger: testLogger,
		Data:   make([]interface{}, 0),
	}

	testFields2 := fields{
		Logger: testLogger2,
		Data:   make([]interface{}, 0),
	}

	testFields3 := fields{
		Logger: testLogger3,
		Data:   make([]interface{}, 0),
	}

	type args struct {
		kvs []interface{}
	}

	testKVS := make([]interface{}, 0)
	testKVS = append(testKVS, "starting", "up and reading for tests")

	testArgs := args{
		kvs: testKVS,
	}

	testKVS2 := make([]interface{}, 0)
	testKVS2 = append(testKVS2, "starting")

	testArgs2 := args{
		kvs: testKVS2,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test-error-write", fields: testFields, args: testArgs},
		{name: "test-error-write", fields: testFields, args: testArgs2},
		{name: "test-error-write", fields: testFields2, args: testArgs},
		{name: "test-error-write", fields: testFields2, args: testArgs2},
		{name: "test-error-write", fields: testFields3, args: testArgs},
		{name: "test-error-write", fields: testFields3, args: testArgs2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &logging.Entry{
				Logger: *tt.fields.Logger,
				Data:   tt.fields.Data,
			}
			e.Info(tt.args.kvs...)
		})
	}
}

// Package network provide some basic network capabilities
package network_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/allens-dev/plumbing/logging"
	"github.com/allens-dev/plumbing/network"
)

func TestHTTPServer(t *testing.T) {
	type args struct {
		parameters *network.ServerParameters
	}

	log := logging.NewJSONLogger("alert-receiver", "")

	port := "9091"

	mux := http.NewServeMux()

	testParameters := &network.ServerParameters{
		Log:  log,
		Port: port,
		Mux:  mux,
	}

	testArgs := args{
		parameters: testParameters,
	}

	testWant := network.HTTPServer(testParameters)

	tests := []struct {
		name string
		args args
		want *network.Server
	}{
		{name: "test-new-server", args: testArgs, want: testWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := network.HTTPServer(tt.args.parameters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

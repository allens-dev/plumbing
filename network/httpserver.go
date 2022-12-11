// Package network provide some basic network capabilities
package network

import (
	"crypto/tls"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"time"

	"github.com/allens-dev/plumbing/logging"
)

const (
	readTimeout       = 30 * time.Second
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 30 * time.Second
)

// ServerParameters provides a way to configure your HTTP Server.
type ServerParameters struct {
	Log          *logging.Logger
	Port         string
	Mux          http.Handler
	Certificates []tls.Certificate
}

// Server provides a HTTP Server with a log attached to it.
type Server struct {
	Log *logging.Logger
	*http.Server
}

// HTTPServer returns a preconfigured HTTP Server.
func HTTPServer(parameters *ServerParameters) *Server {
	return &Server{
		Log: parameters.Log,
		Server: &http.Server{
			Addr:              fmt.Sprintf(":%s", parameters.Port),
			Handler:           parameters.Mux,
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
			TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			TLSConfig: &tls.Config{
				Certificates:     parameters.Certificates,
				MinVersion:       tls.VersionTLS13,
				CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
			ErrorLog: stdLog.New(os.Stderr, "", 0),
		},
	}
}

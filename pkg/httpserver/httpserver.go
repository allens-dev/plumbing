package httpserver

import (
	"crypto/tls"
	"fmt"
	stdLog "log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	readTimeout       = 30 * time.Second
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 30 * time.Second
)

type ServerParameters struct {
	Log          *logrus.Entry
	Port         string
	Mux          http.Handler
	Certificates []tls.Certificate
}

type Server struct {
	Log *logrus.Entry
	*http.Server
}

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
			ErrorLog: stdLog.New(parameters.Log.WriterLevel(logrus.ErrorLevel), "", 0),
		},
	}
}

package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/worker/router"
	"github.com/go-vela/worker/router/middleware"
	"github.com/go-vela/worker/worker"
	"github.com/sirupsen/logrus"
)

// server is a helper function to listen and serve
// traffic for web and API requests for the Worker.
func server(w *worker.Worker) (http.Handler, *tls.Config) {
	// log a message indicating the setup of the server handlers
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Trace
	logrus.Trace("loading router with server handlers")

	// create the worker router to listen and serve traffic
	//
	// https://pkg.go.dev/github.com/go-vela/worker/router?tab=doc#Load
	_server := router.Load(
		WorkerMiddleware(w),
		middleware.RequestVersion,
		middleware.Executors(w.Executors),
		middleware.Secret(w.Config.Server.Secret),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
	)

	// log a message indicating the start of serving traffic
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Tracef
	logrus.Tracef("serving traffic on %s", w.Config.API.Address.Port())

	// if running with HTTPS, check certs are provided and run with TLS.
	if strings.EqualFold(w.Config.API.Address.Scheme, "https") {
		if len(w.Config.Certificate.Cert) > 0 && len(w.Config.Certificate.Key) > 0 {
			// check that the certificate and key are both populated
			_, err := os.Stat(w.Config.Certificate.Cert)
			if err != nil {
				logrus.Fatalf("expecting certificate file at %s, got %v", w.Config.Certificate.Cert, err)
			}

			_, err = os.Stat(w.Config.Certificate.Key)
			if err != nil {
				logrus.Fatalf("expecting certificate key at %s, got %v", w.Config.Certificate.Key, err)
			}
		} else {
			logrus.Fatal("unable to run with TLS: No certificate provided")
		}

		// define TLS config struct for server start up
		tlsCfg := new(tls.Config)

		// if a TLS minimum version is supplied, set that in the config
		if len(w.Config.TLSMinVersion) > 0 {
			var tlsVersion uint16

			switch w.Config.TLSMinVersion {
			case "1.0":
				tlsVersion = tls.VersionTLS10
			case "1.1":
				tlsVersion = tls.VersionTLS11
			case "1.2":
				tlsVersion = tls.VersionTLS12
			case "1.3":
				tlsVersion = tls.VersionTLS13
			default:
				logrus.Fatal("invalid TLS minimum version supplied")
			}

			tlsCfg.MinVersion = tlsVersion
		}

		return _server, tlsCfg
	}

	// else serve over http
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Engine.Run
	return _server, nil
}

// Worker is a middleware function that attaches the
// worker to the context of every http.Request.
func WorkerMiddleware(w *worker.Worker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("worker", *w)
		c.Next()
	}
}

package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// LoggerHttpHandler is a http handler wrapper
func LoggerHttpHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Debugf(
			"%s %s -> %s done in %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	level string

	_error  *log.Logger
	info    *log.Logger
	warning *log.Logger
	debug   *log.Logger

	ERROR  func(string)
	ERRORf func(string, ...interface{})
	INFO   func(string)
	INFOf  func(string, ...interface{})
	WARN   func(string)
	WARNf  func(string, ...interface{})
	DEBUG  func(string)
	DEBUGf func(string, ...interface{})
)

func init() {
	debug = log.New(ioutil.Discard,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	_error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = func(msg string) { _error.Println(msg) }
	ERRORf = func(msg string, v ...interface{}) { _error.Printf(msg, v) }
	SetLogLevel("info")
}

func SetLogLevel(l string) {
	level = l
	switch level {
	case "error":
		INFO = func(string) {}
		INFOf = func(string, ...interface{}) {}
		WARN = func(string) {}
		WARNf = func(string, ...interface{}) {}
		DEBUG = func(string) {}
		DEBUGf = func(string, ...interface{}) {}
	case "info":
		INFO = func(msg string) { info.Println(msg) }
		INFOf = func(msg string, v ...interface{}) { info.Printf(msg, v) }
		WARN = func(string) {}
		WARNf = func(string, ...interface{}) {}
		DEBUG = func(string) {}
		DEBUGf = func(string, ...interface{}) {}
	case "warn":
		INFO = func(msg string) { info.Println(msg) }
		INFOf = func(msg string, v ...interface{}) { info.Printf(msg, v) }
		WARN = func(msg string) { warning.Println(msg) }
		WARNf = func(msg string, v ...interface{}) { warning.Printf(msg, v) }
		DEBUG = func(string) {}
		DEBUGf = func(string, ...interface{}) {}
	case "debug":
		INFO = func(msg string) { info.Println(msg) }
		INFOf = func(msg string, v ...interface{}) { info.Printf(msg, v) }
		WARN = func(msg string) { warning.Println(msg) }
		WARNf = func(msg string, v ...interface{}) { warning.Printf(msg, v) }
		DEBUG = func(msg string) { debug.Println(msg) }
		DEBUGf = func(msg string, v ...interface{}) { debug.Printf(msg, v) }
	default:
		ERRORf("Log level \"%s\" cannot be recognized.", level)
	}
}

func SetErrorOutput(handle io.Writer) {
	_error.SetOutput(handle)
}

func SetInfoOutput(handle io.Writer) {
	info.SetOutput(handle)
}

func SetWarningOutput(handle io.Writer) {
	warning.SetOutput(handle)
}

func SetDebugOutput(handle io.Writer) {
	debug.SetOutput(handle)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

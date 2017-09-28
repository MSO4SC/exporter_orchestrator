package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	monitor = flag.String(
		"monitor-host",
		"",
		"Host and port were monitor (Prometheus) is listening.",
	)
	addr = flag.String(
		"listen-address",
		":8079",
		"The address to listen on for HTTP requests.",
	)
	workDir = flag.String(
		"work-dir",
		"./",
		"Work directory where config.json file is located.",
	)
	logLevel = flag.String(
		"log-level",
		"error",
		"Level of log output.",
	)

	// Global variables of the orchestrator
	config = LoadConfig()
	memory *Memory
)

func main() {
	flag.Parse()

	// Add to workDir ending / if not present
	if (*workDir)[len(*workDir)-1] != '/' {
		*workDir += "/"
	}

	// Parse and set log lovel
	level, err := log.ParseLevel(*logLevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.WarnLevel)
		log.Warnf("Log level %s not recognized, setting \"warn\" as default.")
	}

	// Flags check
	if *monitor == "" {
		flag.Usage()
		log.Fatalln("A host and port must be provided to connect to Prometheus.")
	}

	config = LoadConfig()
	memory = NewMemory()

	if err := memory.LoadFromFile(config.StorageFileName); err != nil {
		log.Warningf("Couldn't load memory file: %s", err.Error())
	}
	memory.StartHealing(30 * time.Second)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(*addr, router))
}

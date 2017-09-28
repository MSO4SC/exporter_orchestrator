package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ExportersIndex handler, returns the memory of the orchestrator
func ExportersIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	memory.Lock()
	defer memory.Unlock()
	if err := memory.Encode(w); err != nil {
		encodeError(w, http.StatusNotFound, err)
		log.Error(err.Error())
		return
	}
}

/*
AddExporter adds an exporter instance to the orchestrator, managing it.
Test with this curl command:

curl -X POST \
  http://localhost:8079/exporters/add \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"host": "ft2.cesga.es",
	"type": "SLURM",
	"persistent": true,
	"args": [
		"user": "[USER]",
		"pass": "[PASS]",
		"tz": "Europe/Madrid",
		"log": "debug"
	]
}'
*/
func AddExporter(w http.ResponseWriter, r *http.Request) {
	modifyExporter(w, r, memory.AddExporterInstance, http.StatusCreated)
}

/*
RemoveExporter removes an exporter instance of the orchestrator, stoping it.
Test with this curl command:

curl -X POST \
  http://localhost:8079/exporters/remove \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"host": "ft2.cesga.es",
	"type": "slurm",
	"persistent": true,
	"args": [
		"-ssh-user [USER]",
		"-ssh-password [PASS]",
		"-countrytz Europe/Madrid",
		"-log.level=warn"
	]
}'
*/
func RemoveExporter(w http.ResponseWriter, r *http.Request) {
	modifyExporter(w, r, memory.RemoveExporterInstance, http.StatusOK)
}

func modifyExporter(w http.ResponseWriter,
	r *http.Request,
	modifier func(*Exporter) error,
	successStatus int) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		encodeError(w, http.StatusNotFound, err)
		log.Error(err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		encodeError(w, http.StatusNotFound, err)
		log.Error(err.Error())
		return
	}

	var exporter Exporter
	if err := json.Unmarshal(body, &exporter); err != nil {
		encodeError(w, 422, err) // StatusUnprocessableEntity not defined in go 1.6
		log.Error(err.Error())
		return
	}

	memory.Lock()
	defer memory.Unlock()
	if err := modifier(&exporter); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		log.Error(err.Error())
		return
	}

	w.WriteHeader(successStatus)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := memory.SaveToFile(config.StorageFileName); err != nil {
		log.Errorf("saving new memory data: %s", err.Error())
	}
}

func encodeError(w http.ResponseWriter, httpCode int, err error) {
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err2 := json.NewEncoder(w).Encode(jsonErr{Code: httpCode, Text: err.Error()}); err2 != nil {
		log.Fatal(err2)
	}
}

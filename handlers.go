package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	//	"github.com/gorilla/mux"
	"time"
)

var (
	memory *Memory
	config = LoadConfig()
)

func init() {
	memory = NewMemory()
	memory.LoadFromFile(config.StorageFileName)
	memory.StartHealing(30 * time.Second)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	memory.Lock()
	defer memory.Unlock()
	if err := memory.Encode(w); err != nil {
		encodeError(w, http.StatusNotFound, err)
		ERROR(err.Error())
		return
	}
}

func ExportersIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	memory.Lock()
	defer memory.Unlock()
	if err := memory.Encode(w); err != nil {
		encodeError(w, http.StatusNotFound, err)
		ERROR(err.Error())
		return
	}
}

/*
AddExporter adds an exporter instance to the orchestrator, managing it.
Test with this curl command:

curl -X POST \
  http://localhost:8080/exporters/add \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: 1474fc98-cf3a-87fb-31f4-8c0c9236da6c' \
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
func AddExporter(w http.ResponseWriter, r *http.Request) {
	modifyExporter(w, r, memory.AddExporterInstance, http.StatusCreated)
}

/*
RemoveExporter removes an exporter instance of the orchestrator, stoping it.
Test with this curl command:

curl -X POST \
  http://localhost:8080/exporters/remove \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: 28ec12dc-772e-db61-e286-a24690fabd80' \
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
		ERROR(err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		encodeError(w, http.StatusNotFound, err)
		ERROR(err.Error())
		return
	}

	var exporter Exporter
	if err := json.Unmarshal(body, &exporter); err != nil {
		encodeError(w, 422, err) // StatusUnprocessableEntity not defined in go 1.6
		ERROR(err.Error())
		return
	}

	memory.Lock()
	defer memory.Unlock()
	if err := modifier(&exporter); err != nil {
		encodeError(w, http.StatusInternalServerError, err)
		ERROR(err.Error())
		return
	}

	w.WriteHeader(successStatus)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	memory.SaveToFile(config.StorageFileName)
}

func encodeError(w http.ResponseWriter, httpCode int, err error) {

	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(jsonErr{Code: httpCode, Text: err.Error()}); err != nil {
		panic(err)
	}
}

// func TodoShow(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	var todoId int
// 	var err error
// 	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
// 		panic(err)
// 	}
// 	todo := RepoFindTodo(todoId)
// 	if todo.Id > 0 {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		if err := json.NewEncoder(w).Encode(todo); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	// If we didn't find it, 404
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusNotFound)
// 	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
// 		panic(err)
// 	}

// }

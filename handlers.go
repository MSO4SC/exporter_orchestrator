package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	//	"github.com/gorilla/mux"
)

const (
	filename = "data.json"
)

var (
	memory *Memory
)

func init() {
	memory = NewMemory()
	memory.LoadFromFile(filename)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	memory.Lock()
	defer memory.Unlock()
	if err := memory.Encode(w); err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
	}
}

func ExportersIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	memory.Lock()
	defer memory.Unlock()
	if err := memory.Encode(w); err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
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
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
	}
	if err := r.Body.Close(); err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
	}
	var exporter Exporter
	if err := json.Unmarshal(body, &exporter); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err) //FIXME(emepetres) handle panic errors with http responses
		}
	}

	memory.Lock()
	defer memory.Unlock()
	if err := memory.AddExporterInstance(&exporter); err == nil {
		w.WriteHeader(http.StatusCreated)
		memory.SaveToFile(filename)
	} else {
		w.WriteHeader(http.StatusInternalServerError) // FIXME(emepetres): support different errors
	}
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
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
	}
	if err := r.Body.Close(); err != nil {
		panic(err) //FIXME(emepetres) handle panic errors with http responses
	}

	var exporter Exporter
	if err := json.Unmarshal(body, &exporter); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err) //FIXME(emepetres) handle panic errors with http responses
		}
	}

	memory.Lock()
	defer memory.Unlock()
	if err := memory.RemoveExporterInstance(&exporter); err == nil {
		w.WriteHeader(http.StatusOK)
		memory.SaveToFile(filename)
	} else {
		w.WriteHeader(http.StatusInternalServerError) // FIXME(emepetres): support different errors
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

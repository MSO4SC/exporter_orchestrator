package main

import (
	"log"
	"net/http"
)

func main() {
	SetLogLevel("debug")

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8079", router))
}

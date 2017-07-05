package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"ExportersIndex",
		"GET",
		"/exporters",
		ExportersIndex,
	},
	Route{
		"AddExporter",
		"POST",
		"/exporters/add",
		AddExporter,
	},
	Route{
		"RemoveExporter",
		"POST",
		"/exporters/remove",
		RemoveExporter,
	},
	// Route{
	// 	"TodoShow",
	// 	"GET",
	// 	"/todos/{todoId}",
	// 	TodoShow,
	// },
}

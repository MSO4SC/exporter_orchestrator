package main

import "net/http"

// Route represents a server route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes holds all the server routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ExportersIndex,
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
}

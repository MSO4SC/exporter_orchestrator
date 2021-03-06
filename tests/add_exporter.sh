#!/bin/bash

curl -X POST \
  http://localhost:8079/exporters/add \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"host": "ft2.cesga.es",
	"type": "SLURM",
	"persistent": false,
	"args": {
		"user": "otarijci",
		"pass": "300tt.yo",
		"tz": "Europe/Madrid",
		"log": "debug"
	}
}'

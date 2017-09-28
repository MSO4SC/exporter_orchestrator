package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Config holds the application configuration
type Config struct {
	StorageFileName       string                       `json:"storageFile"`
	ExportersScripts      map[string]map[string]string `json:"exporters"`
	WaitBeforeHealSeconds int64                        `json:"healWaitSeconds"`
}

// LoadConfig reads the config file parsing its information
func LoadConfig() *Config {
	config := &Config{
		ExportersScripts: make(map[string]map[string]string),
	}
	raw, err := ioutil.ReadFile(*workDir + "config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(raw, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}

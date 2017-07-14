package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	StorageFileName       string                       `json:"storageFile"`
	Monitor               string                       `json:"monitor"`
	ExportersScripts      map[string]map[string]string `json:"exporters"`
	WaitBeforeHealSeconds int64                        `json:"healWaitSeconds"`
}

func LoadConfig() *Config {
	config := &Config{
		ExportersScripts: make(map[string]map[string]string),
	}
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		ERROR(err.Error())
		panic(err)
	}

	err = json.Unmarshal(raw, config)
	if err != nil {
		ERROR(err.Error())
		panic(err)
	}

	return config
}
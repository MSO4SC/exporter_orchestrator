package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Memory stores the different queues of exporters
type Memory struct {
	exporterQueues map[string]*ExporterQueue
	quitHealing    chan struct{}
	healingURL     string
	sync.Mutex
}

// NewMemory creates a new Memory object
func NewMemory() *Memory {
	return &Memory{
		exporterQueues: make(map[string]*ExporterQueue),
		quitHealing:    make(chan struct{}),
		healingURL:     "http://" + *monitor + "/api/v1/query?query=up",
	}
}

// Encode writes the queues as a json in the w writer
func (memo *Memory) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(memo.exporterQueues)
}

// LoadFromFile reads filename and import its data to a Memory object
func (memo *Memory) LoadFromFile(filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// // FIXME(emepetres): Code to restart the Start time of each
	// // queue to current time. Right now is commented because probably
	// // is not needed.
	// err = json.Unmarshal(raw, &memo.exporterQueues)
	// if err == nil {
	// 	for _, queue := range memo.exporterQueues {
	// 		queue.Start = time.Now().Unix()
	// 	}
	// }
	// return err
	return json.Unmarshal(raw, &memo.exporterQueues)
}

// SaveToFile push the Memory queues into the file filename
func (memo *Memory) SaveToFile(filename string) error {
	// Open file and defer close
	var file *os.File
	if _, err := os.Stat("filename"); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("couldn't close memory file: %s", err.Error())
		}
	}()

	// Write file and handle errors
	w := bufio.NewWriter(file)
	if err := memo.Encode(w); err != nil {
		return err
	} else if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

// AddExporterInstance append a new instance of an exporter in its queue
// If the queue does not exists, it creates it and run the new instance
func (memo *Memory) AddExporterInstance(exporter *Exporter) error {
	queue, exists := memo.exporterQueues[exporter.Host]

	//If there is no exporter in execution
	if !exists {
		//New exporter
		memo.exporterQueues[exporter.Host] = NewExporterQueue(exporter)
		return memo.exporterQueues[exporter.Host].Up()
	}

	return queue.Add(exporter)
}

// RemoveExporterInstance deletes an exporter instance from its queue.
func (memo *Memory) RemoveExporterInstance(exporter *Exporter) error {
	queue, exists := memo.exporterQueues[exporter.Host]

	//If there is no exporter in execution
	if !exists {
		return errors.New("Exporter can't be removed, it does not exists")
	}

	err := queue.Remove(exporter)
	if err == nil && !queue.Persistent {
		if queue.Dependencies > 0 {
			err = queue.Up() // If the removed exporter was the current one, we need to run the next in queue
		} else {
			delete(memo.exporterQueues, exporter.Host)
		}
	}
	return err
}

// StartHealing starts the process that execute the heal operation
// over every exporter, every d duration.
// To finish do close(memory.quitHealing)
func (memo *Memory) StartHealing(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ticker.C:
				exporters, err := memo.checkState()
				if err != nil {
					log.Errorf("getting exporters state failed: %s", err.Error())
					continue
				}

				memo.Lock()
				for _, queue := range memo.exporterQueues {
					isUp, exists := exporters[queue.Host]
					err := queue.Heal(exists, isUp)
					if err != nil {
						log.Error("healing: %s", err.Error())
					}
				}
				memo.Unlock()
			case <-memo.quitHealing:
				ticker.Stop()
				return
			}
		}
	}()
}

func (memo *Memory) checkState() (map[string]bool, error) {
	response, err := http.Get(memo.healingURL)
	if err != nil {
		return nil, err
	}

	var decoded map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	result := decoded["data"].(map[string]interface{})["result"].([]interface{})
	exporters := make(map[string]bool)
	for _, entry := range result {
		val, err := strconv.ParseBool(entry.(map[string]interface{})["value"].([]interface{})[1].(string))
		if err != nil {
			log.Warnf("up metric: ", err.Error())
			continue
		}
		exporters[entry.(map[string]interface{})["metric"].(map[string]interface{})["job"].(string)] = val
	}

	return exporters, nil
}

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Memory struct {
	exporterQueues map[string]*ExporterQueue
	quitHealing    chan struct{}
	sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		exporterQueues: make(map[string]*ExporterQueue, 0),
		quitHealing:    make(chan struct{}),
	}
}

func (memo *Memory) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(memo.exporterQueues)
}

func (memo *Memory) LoadFromFile(filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &memo.exporterQueues)
	if err != nil {
		return err
	}

	for _, queue := range memo.exporterQueues {
		if queue.IsUP() {
			queue.Heal()
		}
	}

	return err
}

func (memo *Memory) SaveToFile(filename string) error {
	var file *os.File
	var err error
	if _, err = os.Stat("filename"); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		file, err = os.Open(filename)
	}

	defer file.Close()
	w := bufio.NewWriter(file)
	err = memo.Encode(w)
	w.Flush()
	return err
}

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
				memo.Lock()
				for _, queue := range memo.exporterQueues {
					queue.Heal()
				}
				DEBUG("Healed!")
				memo.Unlock()
			case <-memo.quitHealing:
				ticker.Stop()
				return
			}
		}
	}()
}

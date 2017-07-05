package main

import (
	"errors"
)

type Memory struct {
	exporterQueues map[string]*ExporterQueue
}

func NewMemory() *Memory {
	// TODO(emepetres): Load persistent memory json file
	return &Memory{
		exporterQueues: make(map[string]*ExporterQueue, 0),
	}
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
	// TODO(emepetres): Save persistent memory json file
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
	// TODO(emepetres): Save persistent memory json file
}

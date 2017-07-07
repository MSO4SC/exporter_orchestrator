package main

import (
	"errors"
	"os/exec"
	"reflect"
)

type Exporter struct {
	Host       string   `json:"host"`
	Type       string   `json:"type"`
	Persistent bool     `json:"persistent"`
	Args       []string `json:"args"`
}

func (exporter *Exporter) Create() error {
	cmd := exec.Command(config.ExportersScripts[exporter.Type]["create"],
		exporter.Args...)
	err := cmd.Run()
	ERRORe(err)
	return err
}

func (exporter *Exporter) Heal() error {
	//TODO(emepetres)
	return nil
}

func (exporter *Exporter) Destroy() error {
	cmd := exec.Command(config.ExportersScripts[exporter.Type]["destroy"],
		exporter.Args...)
	err := cmd.Run()
	ERRORe(err)
	return err
}

type ExporterQueue struct {
	Host         string     `json:"host"`
	Type         string     `json:"type"`
	Persistent   bool       `json:"persistent"`
	Dependencies uint       `json:"dep"`
	ArgsQueue    [][]string `json:"queue"`
	Exec         bool       `json:"Exec"`
}

func (exporter *Exporter) belongsToQueue(queue *ExporterQueue) bool {
	if exporter.Host != queue.Host ||
		exporter.Type != queue.Type {
		return false
	}
	return true
}

func NewExporterQueue(exp *Exporter) *ExporterQueue {
	return &ExporterQueue{
		Host:         exp.Host,
		Type:         exp.Type,
		Persistent:   exp.Persistent,
		Dependencies: 1,
		ArgsQueue:    [][]string{exp.Args},
		Exec:         false,
	}
}

func (expQ *ExporterQueue) IsUP() bool {
	return expQ.Exec
}

func (expQ *ExporterQueue) Up() error {
	if expQ.Exec {
		return nil
	}
	err := expQ.getCurrentExporter().Create()
	expQ.Exec = (err == nil)
	return err
}

func (expQ *ExporterQueue) Down() error {
	if !expQ.Exec {
		return nil
	}
	err := expQ.getCurrentExporter().Destroy()
	expQ.Exec = (err != nil)
	return err
}

func (expQ *ExporterQueue) Heal() error {
	if !expQ.Exec || (expQ.Dependencies == 0 && !expQ.Persistent) {
		return nil
	}
	err := expQ.getCurrentExporter().Heal()
	expQ.Exec = (err == nil)
	return err
}

// Add adds a new exporter to the queue.
func (expQ *ExporterQueue) Add(exp *Exporter) error {
	if !exp.belongsToQueue(expQ) {
		return errors.New("exporter with host " + exp.Host + " does not belongs to queue")
	}

	if expQ.Persistent {
		expQ.Dependencies++
		return nil
	}

	//TODO(emepetres) What happens if new exporter is persistent?

	var err error
	expQ.ArgsQueue = append(expQ.ArgsQueue, exp.Args)
	expQ.Dependencies++

	return err
}

// Remove deletes an exporter in the queue.
// If the exporter is the current one, it stops it before deleting.
func (expQ *ExporterQueue) Remove(exp *Exporter) error {
	if expQ.Dependencies == 0 {
		return errors.New("trying to change exporter %s with no dependencies left")
	}

	if !exp.belongsToQueue(expQ) {
		return errors.New("Exporter with host " + exp.Host + " does not belongs to queue.")
	}

	if expQ.Persistent {
		expQ.Dependencies--
		return nil
	}

	// Get exporter index
	i := expQ.findExporter(exp)
	if i < 0 {
		return errors.New("cannot remove exporter, it doesn't exists in the queue")
	}

	// Remove current running instance on a non persistent exporter
	if i == 0 && expQ.IsUP() {
		err := expQ.Down()
		if err == nil {
			expQ.Dependencies--
			if expQ.Dependencies == 0 {
				expQ.ArgsQueue = make([][]string, 0)
			} else {
				expQ.ArgsQueue = expQ.ArgsQueue[1:]
			}
		}
		return err
	}

	// Remove not running instance on a non persistent exporter
	expQ.Dependencies--
	if expQ.Dependencies == 0 {
		expQ.ArgsQueue = make([][]string, 0)
	} else {
		if i < (len(expQ.ArgsQueue) - 1) {
			expQ.ArgsQueue = append(expQ.ArgsQueue[:i], expQ.ArgsQueue[i+1:]...)
		} else {
			expQ.ArgsQueue = expQ.ArgsQueue[:i]
		}
	}
	return nil

}

func (expQ *ExporterQueue) getCurrentExporter() *Exporter {
	return &Exporter{
		Host:       expQ.Host,
		Type:       expQ.Type,
		Persistent: expQ.Persistent,
		Args:       expQ.ArgsQueue[0],
	}
}

func (expQ *ExporterQueue) findExporter(exp *Exporter) int {
	for i, args := range expQ.ArgsQueue {
		if reflect.DeepEqual(args, exp.Args) {
			return i
		}
	}
	return -1
}

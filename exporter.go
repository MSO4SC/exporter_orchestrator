package main

import (
	"errors"
	"net"
	"os/exec"
	"reflect"
	"time"
)

type Exporter struct {
	Host       string            `json:"host"`
	Type       string            `json:"type"`
	Persistent bool              `json:"persistent"`
	Args       map[string]string `json:"args"`
}

func (exporter *Exporter) Create(listenPort string) error {
	cmd := exec.Command(config.ExportersScripts[exporter.Type]["create"],
		listenPort,
		exporter.Host,
		exporter.Args["user"],
		exporter.Args["pass"],
		exporter.Args["tz"],
		exporter.Args["log"])
	err := cmd.Run()
	if err != nil {
		ERROR(err.Error())
	}
	return err
}

func (exporter *Exporter) Destroy(listenPort string) error {
	cmd := exec.Command(config.ExportersScripts[exporter.Type]["destroy"],
		listenPort,
		exporter.Host,
		exporter.Args["user"],
		exporter.Args["pass"],
		exporter.Args["tz"],
		exporter.Args["log"])
	err := cmd.Run()
	if err != nil {
		ERROR(err.Error())
	}
	return err
}

func (exporter *Exporter) belongsToQueue(queue *ExporterQueue) bool {
	if exporter.Host != queue.Host ||
		exporter.Type != queue.Type {
		return false
	}
	return true
}

type ExporterQueue struct {
	Host         string              `json:"host"`
	Type         string              `json:"type"`
	ListenPort   string              `json:"listen-port"`
	Persistent   bool                `json:"persistent"`
	Dependencies uint                `json:"dep"`
	ArgsQueue    []map[string]string `json:"queue"`
	Exec         bool                `json:"Exec"`
	Start        int64               `json:"start"`
}

func NewExporterQueue(exp *Exporter) *ExporterQueue {
	return &ExporterQueue{
		Host:         exp.Host,
		Type:         exp.Type,
		ListenPort:   getFreePort(),
		Persistent:   exp.Persistent,
		Dependencies: 1,
		ArgsQueue:    []map[string]string{exp.Args},
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
	err := expQ.getCurrentExporter().Create(expQ.ListenPort)
	expQ.Exec = (err == nil)
	if expQ.IsUP() {
		expQ.Start = time.Now().Unix()
	}
	return err
}

func (expQ *ExporterQueue) Down() error {
	if !expQ.Exec {
		return nil
	}
	err := expQ.getCurrentExporter().Destroy(expQ.ListenPort)
	expQ.Exec = (err != nil)
	return err
}

func (expQ *ExporterQueue) Heal(exists, isUp bool) error {
	if time.Now().Unix()-expQ.Start < config.WaitBeforeHealSeconds {
		return nil
	}

	if !expQ.IsUP() || !exists || !isUp {
		expQ.Exec = false
		WARN("healing " + expQ.Host + "...")
		return expQ.Up()
	}

	return nil
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
		return errors.New("trying to change exporter " + expQ.Host + " with no dependencies left")
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
				expQ.ArgsQueue = make([]map[string]string, 0)
			} else {
				expQ.ArgsQueue = expQ.ArgsQueue[1:]
			}
		}
		return err
	}

	// Remove not running instance on a non persistent exporter
	expQ.Dependencies--
	if expQ.Dependencies == 0 {
		expQ.ArgsQueue = make([]map[string]string, 0)
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
	// reverse search to find first the ones not executing
	for i := len(expQ.ArgsQueue) - 1; i >= 0; i-- {
		if reflect.DeepEqual(expQ.ArgsQueue[i], exp.Args) {
			return i
		}
	}
	return -1
}

func getFreePort() string {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	return l.Addr().String()[4:]
}

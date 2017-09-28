# exporter_orchestrator

Orchestrates the creation, healing, destruction and discovery of exporters through an HTTP API.

## Install

> Requires Go >=1.8

```
go get github.com/mso4sc/slurm_exporter
$GOPATH/src/github.com/mso4sc/slurm_exporter/utils/install.sh
```

## Usage

```bash
slurm_exporter -monitor-host=<HOST:MPORT> [-listen-address=:<PORT>] [-log.level=<LOGLEVEL>]
```

### Defaults

PORT: `:8079`  
LOGLEVEL: `error`  

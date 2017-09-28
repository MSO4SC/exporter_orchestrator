# MSO4SC exporter orchestrator Docker image

Some scripts helps working with docker:  
`run.sh` runs a new exporter orchestrator in a new container. It returns the container ID.  
`stop.sh` stops the container passed as parameter (by ID or NAME).  

### Usage

```
$ ./run.sh  -monitor-host=localhost:9090
ea994b6b6ac2c73f10ca2a1150e32938031ad98a786dab5554772140c1a35c16

$ docker ps -a
CONTAINER ID        IMAGE                   COMMAND                  CREATED         STATUS              PORTS                     NAMES
ea994b6b6ac2        mso4sc/exporter_orc...  "slurm_exporter -l..."   7 minutes ago   Up 3 seconds        0.0.0.0:8079->8079/tcp   dreamy_spence

$ ./stop.sh ea994b6b6ac2c73f10ca2a1150e32938031ad98a786dab5554772140c1a35c16
ea994b6b6ac2c73f10ca2a1150e32938031ad98a786dab5554772140c1a35c16
```

## Development

Two scripts help building and publishing the image  
`build.sh` build the image using the Dockerfile  
`publish.sh` push the image in Docker Hub  

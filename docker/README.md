# MSO4SC exporter orchestrator Docker image

The exporter orchestrator will try to create and orchestrate exporters by its HTTP API. This is done through config and scripts on /app folder. To change its default configuration load /app as a volume and add your own config and script files. 

### Usage
Before starting the exporter orchestrator, an MSOMonitor needs to be running and reachable. In the MSO4SC GitHub examples you can find a docker-compose example to load both of them.

The default configuration uses HOST docker to create new exporters. In order to do that, host's docker has to be mounted on the guest:
```
# docker run --rm -d -p 8079:8079 \
        -v /lib64:/lib64 -v /usr:/usr -v /lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu -v /var/run/docker.sock:/var/run/docker.sock \
        mso4sc/exporter_orchestrator -monitor-host=localhost:9090
ea994b6b6ac2c73f10ca2a1150e32938031ad98a786dab5554772140c1a35c16

# docker ps -a
CONTAINER ID        IMAGE                   COMMAND                  CREATED         STATUS              PORTS                     NAMES
ea994b6b6ac2        mso4sc/exporter_orc...  "exporter_orchestr..."   7 minutes ago   Up 3 seconds        0.0.0.0:8079->8079/tcp   dreamy_spence
```

One script in docker folder helps running the image:  
`run.sh` runs a new exporter orchestrator in a new container. It returns the container ID.

## Development
The image has auto-build in DockerHub. Nevertheless, two scripts help building and publishing the image. 
`build.sh` build the image using the Dockerfile  
`publish.sh` push the image in Docker Hub  

# tcp-wait

Go package to test and wait on the availability of a TCP host and port.
This package is Go port of [wait-for-it.sh](https://github.com/vishnubob/wait-for-it)

## Building

First clone down the repository then use the make file to build what you need. There is a simple single binary build below, but you also have `build-all` and other commands available in the makefile.

```bash
# build and run using make all
$ make all
go get
go test -v ./...
=== RUN   TestSuccessSingle
--- PASS: TestSuccessSingle (0.01s)
=== RUN   TestFailureDouble
{"level":"warning","msg":"tcp ping failed","tcp-host":"nowhere:50","time":"2020-03-13T11:38:45+10:30"}
{"level":"warning","msg":"tcp ping failed","tcp-host":"nowhere:51","time":"2020-03-13T11:38:45+10:30"}
--- PASS: TestFailureDouble (0.50s)
PASS
ok      tcp-wait        (cached)
go build -o ./bin/tcp-wait -v
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/tcp-wait -v
tcp-wait
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/mac/tcp-wait -v
tcp-wait

# your os binary should be imediately executable
$ ./bin/tcp-wait
Usage of ./bin/tcp-wait:
  -it value
        <host:port> [host2:port,...] comma seperated list of services
  -o string
        output in format json/text (default "json")
  -t int
        timeout (default 20)
```


## Usage

```bash
### simple
$ ./bin/tcp-wait -it github.com:80
{"level":"info","msg":"services are ready!","services":["github.com:80"],"time":"2020-03-12T17:18:30+10:30"}

### multiple hosts with timeout and text
$ ./bin/tcp-wait -it github.com:443,google.com:443 -t 1 -o text
INFO[2020-03-12T17:20:15+10:30] services are ready!  services="[github.com:443 google.com:443]"

### multiple hosts with fail condition
$ ./bin/tcp-wait -it github.com:443,localhost:10000 -t 2
{"level":"warning","msg":"tcp ping failed","tcp-host":"localhost:10000","time":"2020-03-12T17:26:16+10:30"}
{"level":"warning","msg":"tcp ping failed","tcp-host":"localhost:10000","time":"2020-03-12T17:26:17+10:30"}
{"level":"error","msg":"services did not respond","time":"2020-03-12T17:26:18+10:30"}
```


<!-- ### docker users
~2.7MB docker image.
```bash
$ docker run --rm alioygur/wait-for -it google.com:80
services are ready!
``` -->


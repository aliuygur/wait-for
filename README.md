# tcp-wait

Go package to test and wait on the availability of a TCP host and port. This is mainly used in some
containers as a pre stat script to wait for ports to be open or die if they don't. It has come in
handy for other tasks as well like watching ports while you are opening and closing.

Feel free to contribute as I will allow most merges as long as they don't drasticly break the usage.

## Building

First clone down the repository then use the make file to build what you need. There is a simple
 single binary build below, but you also have `build-all` and other commands available in the makefile.

## Running

You can install this using go get and adding to your gopath. You can then use as you would any packages.
If you want to use without needing the full path you can add your go path to your shell.

```bash
# get the package
go get github.com/donkeyx/tcp-wait

# add go bin path to your startup shell, bash in this case
echo 'export PATH="~/go/bin:$PATH"' ~/.bashrc

# quick run from there with
tcp-wait -hp localhost:8080 -t 5

# for help just run with no flags
tcp-wait -h
Usage of /Users/dbinney/go/bin/tcp-wait:
  -hp value
    	<host:port> [host2:port,...] comma seperated list of host:ports
  -o string
    	output in format json/text (default "json")
  -t int
    	timeout (default 20)
  -version
    	version information
...
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

## Building locally

You can just clone down the repo, then use the make file to build the packages. They will be placed
in a local bin folder, with linux/mac and a local env build. You can execute them directly from that
path.

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

# structure created with binaries in your local bin directory.
$ tree
.
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── bin
│   ├── linux
│   │   └── tcp-wait
│   ├── mac
│   │   └── tcp-wait
│   └── tcp-wait
├── go.mod
├── go.sum
├── main.go
└── main_test.go


# Binary should be imediately executable
$ ./bin/tcp-wait
Usage of ./bin/tcp-wait:
  -it value
        <host:port> [host2:port,...] comma seperated list of services
  -o string
        output in format json/text (default "json")
  -t int
        timeout (default 20)
```


<!-- ### docker users
~2.7MB docker image.
```bash
$ docker run --rm alioygur/wait-for -it google.com:80
services are ready!
``` -->


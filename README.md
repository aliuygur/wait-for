# wait-for

Go package to test and wait on the availability of a TCP host and port.
This package is Go port of [wait-for-it.sh](https://github.com/vishnubob/wait-for-it)

# Usage

```bash
wait-for:
  -it value
        <host:port> [host2:port,...] comma seperated list of services
  -t int
        timeout (default 20)
```

# Example


### simple
```bash
$ wait-for -it github.com:80 && echo "github is up!"
services are ready!
github is up!
```

### multiple hosts and custom timeout
```bash
$ wait-for -t 5 -it github.com:80,bitbucket.com:80 && echo "github and bitbucket are up!"
services are ready!
github and bitbucket are up!
```

### docker users
~2.7MB docker image.
```bash
$ docker run --rm alioygur/wait-for -it google.com:80
services are ready!
```

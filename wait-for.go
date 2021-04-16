package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type servicesType []string

var timeout int
var services servicesType

func (s *servicesType) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *servicesType) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// waitForServices tests and waits on the availability of a TCP host and port
func waitForServices(services []string, timeOut time.Duration) error {
	var depChan = make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(services))
	go func() {
		for _, s := range services {
			go func(s string) {
				defer wg.Done()
				for {
					_, err := net.Dial("tcp", s)
					if err == nil {
						return
					}
					time.Sleep(1 * time.Second)
				}
			}(s)
		}
		wg.Wait()
		close(depChan)
	}()

	select {
	case <-depChan: // services are ready
		return nil
	case <-time.After(timeOut):
		return fmt.Errorf("services aren't ready in %s", timeOut)
	}
}

func init() {
	flag.IntVar(&timeout, "t", 20, "timeout")
	flag.Var(&services, "it", "<host:port> [host2:port,...] comma seperated list of services")
}

func main() {
	flag.Parse()
	if len(services) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := waitForServices(services, time.Duration(timeout)*time.Second); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("services are ready!")
	os.Exit(0)
}

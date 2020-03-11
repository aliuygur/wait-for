package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
					log.WithFields(log.Fields{
						"tcp-host": s,
					}).Warn("tcp ping failed")
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
		log.WithFields(log.Fields{
			"timeout": timeOut,
		}).Error("services did not repont in time")
		return fmt.Errorf("")
	}
}

func init() {

	// log.SetFormatter(&log.TextFormatter{
	// 	DisableColors: false,
	// 	FullTimestamp: true,
	// })

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)
	log.SetOutput(os.Stdout)
	// log = log.WithFields(log.Fields{"service": "tcp-wait"})

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
		os.Exit(1)
	}
	log.Info("services are ready!")
	os.Exit(0)
}

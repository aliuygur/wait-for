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

//version info set with ldflags
var (
	Version string = "undefined"
	GitHash string = "undefined"
	// BuildDate      string = "date"
	displayVersion bool
)

func (s *servicesType) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *servicesType) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

func GetVersion() {
	fmt.Printf("Version %s\n", Version)
	fmt.Printf("Git_Hash %s\n", GitHash)
	// fmt.Printf("Build_Date %s\n", BuildDate)
}

// waitForServices tests and waits on the availability of a TCP host and port
func waitForServices(services []string, timeOut time.Duration) error {
	// fmt.Println(services, len(services), timeOut)
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
		return fmt.Errorf("services did not respond")
	}
}

func init() {

	flag.BoolVar(&displayVersion, "version", false, "version information")
	flag.IntVar(&timeout, "t", 20, "timeout")
	flag.Var(&services, "hp", "<host:port> [host2:port,...] comma seperated list of host:ports")

	output := flag.String("o", "json", "output in format json/text")

	if *output == "text" {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	flag.Parse()

	// exit on version request
	if displayVersion {
		GetVersion()
		os.Exit(0)
	}
}

func main() {

	if len(services) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	resp := waitForServices(services, time.Duration(timeout)*time.Second)
	if resp != nil {
		log.Error(resp)
		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"services": services,
	}).Info("services are ready!")
	// log.Infof("services are ready! %s", services)
	os.Exit(0)

}

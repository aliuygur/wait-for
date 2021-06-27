package wfi

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type check struct {
	tcp  *tcpCheckService
	http *url.URL
}

func (c check) String() string {
	if c.tcp != nil {
		return c.tcp.String()
	}
	if c.http != nil {
		return c.http.String()
	}
	return ""
}

func (c check) execute() error {
	if c.tcp != nil {
		return c.executeTcp()
	}
	if c.http != nil {
		return c.executeHttp()
	}
	return fmt.Errorf("Invalid check, has nothing to check")
}

func (c check) executeTcp() error {
	checkHost := c.tcp.String()

	log.Debugf("TCP [%s]", checkHost)

	_, err := net.Dial("tcp", c.tcp.String())

	return err
}

func (c check) executeHttp() error {
	checkUrl := c.http.String()

	log.Debugf("HTTP [GET %s]", checkUrl)

	_, err := http.Get(checkUrl)

	return err
}

type tcpCheckService struct {
	host string
	port int
}

func (c tcpCheckService) String() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

type WaitForChecker interface {
	Execute() error
}

const DefaultWaitForBeforeStart = 0
const DefaultSleepBetweenChecks = 1
const DefaultTimeoutInSeconds = 30

type service struct {
	checks             []check
	waitBeforeStart    int
	sleepBetweenChecks int
	timeout            int
}

func New(services []string, wait int, timeout int, sleepBetweenChecks int) (WaitForChecker, error) {
	svc := &service{
		waitBeforeStart:    wait,
		timeout:            timeout,
		sleepBetweenChecks: sleepBetweenChecks,
		checks:             []check{},
	}

	if svc.waitBeforeStart < 1 {
		svc.waitBeforeStart = DefaultWaitForBeforeStart
	}

	if svc.timeout < 1 {
		svc.timeout = DefaultTimeoutInSeconds
	}

	if svc.sleepBetweenChecks < 1 {
		svc.sleepBetweenChecks = DefaultSleepBetweenChecks
	}

	var err error

	tcpServices, err := servicesToTcpChecks(services)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot get tcp checks")
	}
	for _, c := range tcpServices {
		svc.checks = append(svc.checks, c)
	}

	httpServices, err := servicesToHttpChecks(services)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot get http checks")
	}
	for _, c := range httpServices {
		svc.checks = append(svc.checks, c)
	}

	return svc, nil
}

func servicesToTcpChecks(rawServices []string) ([]check, error) {
	const prefix = "tcp://"

	checks := []check{}

	for _, rawService := range rawServices {
		if strings.HasPrefix(rawService, prefix) {
			rawValue := strings.TrimPrefix(rawService, prefix)
			if strings.Contains(rawValue, ":") {
				parts := strings.Split(rawValue, ":")
				if len(parts) == 2 {
					host := parts[0]
					port, err := strconv.Atoi(parts[1])
					if err != nil {
						return nil, errors.Wrapf(err, "Invalid tcp value [%s]", rawService)
					}

					checks = append(checks, check{
						tcp: &tcpCheckService{
							host: host,
							port: port,
						},
					})

				} else {
					return nil, fmt.Errorf("Invalid tcp value [%s]", rawService)
				}
			} else {
				return nil, fmt.Errorf("Invalid tcp value [%s]", rawService)
			}
		}
	}

	return checks, nil
}

func servicesToHttpChecks(rawServices []string) ([]check, error) {
	checks := []check{}

	for _, rawService := range rawServices {
		if strings.HasPrefix(rawService, "http") || strings.HasPrefix(rawService, "https://") {
			serviceUrl, err := url.Parse(rawService)
			if err != nil {
				return nil, errors.Wrapf(err, "Cannot convert [%s] to URL", rawService)
			}

			checks = append(checks, check{
				http: serviceUrl,
			})
		}
	}

	return checks, nil
}

func (svc *service) Execute() error {
	if svc.waitBeforeStart > 0 {
		waitInSeconds := time.Duration(svc.waitBeforeStart) * time.Second

		log.Tracef("Wait in %v", waitInSeconds)

		time.Sleep(waitInSeconds)
	}

	if err := svc.waitForServices(time.Duration(svc.timeout)*time.Second, time.Duration(svc.sleepBetweenChecks)*time.Second); err != nil {
		return err
	}

	return nil
}

func (svc *service) waitForServices(timeout time.Duration, sleepBetwenChecks time.Duration) error {
	log.Tracef("Start checking with timeout %v", timeout)

	var depChan = make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(svc.checks))
	go func() {
		for _, c := range svc.checks {
			go func(c check) {
				defer wg.Done()
				for {
					err := c.execute()
					if err == nil {
						log.Tracef("Check [%s] OK", c.String())

						return
					}

					log.Debugln(errors.Wrapf(err, "Cannot check [%s]", c.String()))
					log.Tracef("Wait in %v", sleepBetwenChecks)

					time.Sleep(sleepBetwenChecks)
				}
			}(c)
		}

		wg.Wait()
		close(depChan)
	}()

	select {
	case <-depChan: // services are ready
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("services aren't ready in %s", timeout)
	}
}

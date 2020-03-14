package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestSuccessSingle(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	var duration time.Duration = 500 * 1000000 // 100ms

	hosts := []string{"google.com:80"}

	var resp = waitForServices(hosts, duration)
	assert.Equal(t, resp, nil, "Success response should be nil")
}

func TestFailureDouble(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	var duration time.Duration = 500 * 1000000 // 100ms

	hosts := []string{"nowhere:50", "nowhere:51"}

	var resp = waitForServices(hosts, duration)
	err := errors.New("services did not respond")

	assert.Equal(t, resp, err, "Error message does not match.")
}

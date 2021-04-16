IMAGE_NAME ?= owncloudci/wait-for:latest

BUILD_VERSION ?= latest
BUILD_ARCH ?= amd64
BUILD_DATE ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

VCS_REF ?= $(shell git rev-parse HEAD)
VCS_URL ?= $(shell git remote get-url origin)

DOCKER_CMD ?= docker

.PHONY: build
build:
	cd $(BUILD_VERSION) && $(DOCKER_CMD) build -f Dockerfile.$(BUILD_ARCH) --label org.label-schema.version=$(BUILD_VERSION) --label org.label-schema.build-date=$(BUILD_DATE) --label org.label-schema.vcs-url=$(VCS_URL) --label org.label-schema.vcs-ref=$(VCS_REF) -t $(IMAGE_NAME) .

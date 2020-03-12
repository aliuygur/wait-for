
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_PATH=./bin/linux/
BINARY_NAME=tcp-wait
BINARY_UNIX=$(BINARY_NAME)

all: test build
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	find ./bin/ -type f | grep -v keep | xargs rm
# run:
# 	$(GOBUILD) -o $(BINARY_NAME) -v ./...
# 	./$(BINARY_NAME)
deps:
	$(GOGET)


# Cross compilation
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/linux/$(BINARY_NAME) -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/mac/$(BINARY_NAME) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

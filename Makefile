
# Go parameters
BINARY_PATH = ./bin/linux/
BINARY_NAME = tcp-wait
VERSION ='$(shell git describe --tags)'
VERSION ='$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)'
BUILD_DATE='$(shell date)'
HASH = '$(shell git rev-parse --short HEAD)'
BUILD_FLAGS = go build -ldflags "-w -s -X main.Version=$(VERSION) -X main.GitHash=$(HASH)"


all: clean deps test build-all
build:
	go build -o ./bin/$(BINARY_NAME) -v
test:
	mkdir -p tmp/test-coverage
	go test -coverprofile=tmp/test-coverage/coverage.out
	go tool cover -html=tmp/test-coverage/coverage.out -o ./tmp/test-coverage/coverage.html
clean:
	go clean
	find ./bin/ -type f | grep -v keep | xargs rm
# run:
# 	$(GOBUILD) -o $(BINARY_NAME) -v ./...
# 	./$(BINARY_NAME)
deps:
	go get


# Cross compilation
build-all:
	$(info    version is $(VERSION))
	$(info    build_date is $(BUILD_DATE))
	$(info    ld-flags is $(BUILD_FLAGS))

	$(BUILD_FLAGS) -o ./bin/$(BINARY_NAME) -v
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_FLAGS) -o bin/linux/$(BINARY_NAME) -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD_FLAGS) -o bin/mac/$(BINARY_NAME) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_NAME)" -v

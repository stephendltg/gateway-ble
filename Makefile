# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOLINT=golint
BINARY_NAME=gateway-ble

all: deps tool build

dev:
	$(GORUN) main.go -du=10s -mqtt=127.0.0.1:1883 -db=http://127.0.0.1:8086 -debug

build:
	$(GOBUILD) -v .

build-linux:
	GOOS=linux $(GOBUILD) -v -o $(BINARY_NAME)-linux .

build-rasp:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o $(BINARY_NAME)-rasp .

tool:
	$(GOVET) ./...; true
	$(GOFMT) -w .

clean:
	go clean -i .
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME)-rasp

deps:
	go mod tidy
	go mod verify

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down -v

docker-build: clean build
	docker build -t $(BINARY_NAME) .
 	docker run --rm $(BINARY_NAME):latest

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make clean: remove object files and cached files"
	@echo "make deps: get the deployment tools"
	@echo "make docker-run: Start docker"
	@echo "make docker-stop: get Stop docker"

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
AUTHOR=stephendltg

all: deps tool build-linux build-rasp

dev:
	$(GORUN) main.go -du=0 -mqtt=127.0.0.1:1883 -db=http://127.0.0.1:8086 -debug

build-app:
	$(GOBUILD) -v -race .

build-linux:
	GOOS=linux $(GOBUILD) -v -o build/$(BINARY_NAME)-linux .

build-rasp:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o build/$(BINARY_NAME)-rasp .

tool:
	$(GOVET) ./...; true
	$(GOFMT) -w .

clean:
	go clean -i .
	rm -f $(BINARY_NAME)
	rm -f build/$(BINARY_NAME)-linux
	rm -f build/$(BINARY_NAME)-rasp
	rm -f *.db

deps:
	go mod vendor
	go mod verify

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down -v

docker-build: clean build
	docker build -t $(AUTHOR)/$(BINARY_NAME) .
 	docker run --rm $(AUTHOR)/$(BINARY_NAME):latest

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make clean: remove object files and cached files"
	@echo "make deps: get the deployment tools"
	@echo "make docker-run: Start docker"
	@echo "make docker-stop: get Stop docker"

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
VERSION := $(shell node -p "require('./package.json').version")
HOMEPAGE := $(shell node -p "require('./package.json').homepage")
DESCRIPTION := $(shell node -p "require('./package.json').description")
PKG_LINUX=build/gateway-linux
PKG_RASP=build/gateway-rasp
NODE=v14.16.1
NVM=v0.38.0

all: deps tool build-linux build-rasp build-deb build-deb-rasp

pre-install: 
	@echo "Installing project ${BINARY_NAME}..."
	. ${NVM_DIR}/nvm.sh && nvm install ${NODE} && nvm use ${NODE}

dev:
	$(GORUN) main.go -du=0 -mqtt=127.0.0.1:1883 -debug -db=http://127.0.0.1:8086 -collect

other:
	$(GORUN) main.go -du=0 -mqtt=127.0.0.1:1883 -debug -mac=ac:23:3f:58:a9:d1 -collect -db=http://127.0.0.1:8086

build-app:
	$(GOBUILD) -v -race .

build-linux:
	GOOS=linux $(GOBUILD) -v -o build/$(BINARY_NAME)-linux .

build-rasp:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o build/$(BINARY_NAME)-rasp .

build-deb:
	mkdir -p $(PKG_LINUX)/DEBIAN
	mkdir -p $(PKG_LINUX)/usr/bin/
	echo "Package: $(BINARY_NAME)" > $(PKG_LINUX)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(PKG_LINUX)/DEBIAN/control
	echo "Section: custom" >> $(PKG_LINUX)/DEBIAN/control
	echo "Architecture: all" >> $(PKG_LINUX)/DEBIAN/control
	echo "Essential: no" >> $(PKG_LINUX)/DEBIAN/control
	echo "Maintainer: $(AUTHOR)" >> $(PKG_LINUX)/DEBIAN/control
	echo "Description: $(DESCRIPTION)" >> $(PKG_LINUX)/DEBIAN/control
	echo "Homepage: $(HOMEPAGE)" >> $(PKG_LINUX)/DEBIAN/control
	GOOS=linux $(GOBUILD) -v -o $(PKG_LINUX)/usr/bin/$(BINARY_NAME) .
	sudo dpkg-deb --build $(PKG_LINUX)
	rm -r $(PKG_LINUX)/*
	rmdir $(PKG_LINUX)

build-deb-rasp:
	mkdir -p $(PKG_RASP)/DEBIAN
	mkdir -p $(PKG_RASP)/usr/bin/
	echo "Package: $(BINARY_NAME)" > $(PKG_RASP)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(PKG_RASP)/DEBIAN/control
	echo "Section: custom" >> $(PKG_RASP)/DEBIAN/control
	echo "Architecture: all" >> $(PKG_RASP)/DEBIAN/control
	echo "Essential: no" >> $(PKG_RASP)/DEBIAN/control
	echo "Maintainer: $(AUTHOR)" >> $(PKG_RASP)/DEBIAN/control
	echo "Description: $(DESCRIPTION)" >> $(PKG_RASP)/DEBIAN/control
	echo "Homepage: $(HOMEPAGE)" >> $(PKG_RASP)/DEBIAN/control
	GOOS=linux $(GOBUILD) -v -o $(PKG_RASP)/usr/bin/$(BINARY_NAME) .
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o $(PKG_RASP)/usr/bin/$(BINARY_NAME) .
	sudo dpkg-deb --build $(PKG_RASP)
	rm -r $(PKG_RASP)/*
	rmdir $(PKG_RASP)

tool:
	$(GOVET) ./...; true
	$(GOFMT) -w .

clean:
	go clean -i .
	rm -f $(BINARY_NAME)
	rm -f build/$(BINARY_NAME)-linux
	rm -f build/$(BINARY_NAME)-rasp
	rm -f build/*.deb
	rm -f *.db

deps:
	# go mod tidy
	go mod vendor
	go mod verify

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down -v

docker-build: clean build
	docker build -t $(AUTHOR)/$(BINARY_NAME) .
 	docker run --rm $(AUTHOR)/$(BINARY_NAME):latest

nvm:
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/${NVM}/install.sh | bash

help:
	@echo "make deps: get the deployment tools"
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make clean: remove object files and cached files"
	@echo "make nvm: insall nvm"
	@echo "make pre-install: Pre install nodejs"
	@echo "make docker-run: Start docker"
	@echo "make docker-stop: get Stop docker"

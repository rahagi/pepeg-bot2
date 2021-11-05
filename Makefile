OUT := pepeg-bot2
PKG := github.com/rahagi/pepeg-bot2
VERSION := $(shell git describe --always --long) 
DOCKER_TAG := ${OUT}:latest
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

all: run

image:
	docker build -t ${DOCKER_TAG} --build-arg VERSION=${VERSION} .

compose:
	docker-compose up

composed:
	docker-compose up -d

test:
	@go test -short ${PKG_LIST}

clean:
	docker image rm -f ${DOCKER_TAG}

run: image compose

.PHONY: all run
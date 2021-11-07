OUT := pepeg-bot2
PKG := github.com/rahagi/pepeg-bot2
TAG := $(shell git tag)
SHA8 := $(shell git rev-parse --short HEAD)
VERSION := ${TAG}-${SHA8}
DOCKER_TAG := ${OUT}:latest
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

all: run

image:
	docker build -t ${DOCKER_TAG} --build-arg VERSION=${VERSION} .

compose:
	docker-compose down && docker-compose up --remove-orphans -d

train:
	docker-compose -f docker-compose.train.yml up -d

test:
	@go test -short ${PKG_LIST}

clean:
	docker image rm -f ${DOCKER_TAG} && docker-compose down --remove-orphans

run: image compose

.PHONY: all run
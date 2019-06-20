#!/usr/bin/env bash

SHELL = /bin/bash
all: clean test

clean:
	@docker-compose rm -f postgresql

test:
	@docker run -d --name postgres  -e ALLOW_EMPTY_PASSWORD=yes circleci/postgres:latest
	@docker run --name golang --rm -v $(PWD):/usr/src/myapp -w /usr/src/myapp  -e TEST_DB_HOST=postgresql --link=postgres:postgresql golang:1.12 go test -v
	@docker rm -f postgres

build:
	@docker run --rm -v $(PWD):/usr/src/myapp -w /usr/src/myapp  golang:1.12 go build

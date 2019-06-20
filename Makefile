#!/usr/bin/env bash

SHELL = /bin/bash
all: clean test

clean:
	@docker-compose rm -f postgresql

test:
	@go get ./...
	@docker-compose up -d postgresql
	@go test -v
	@docker-compose rm -f postgresql

build:
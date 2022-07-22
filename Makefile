MAKEFLAGS += --silent

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

export GONOSUMDB=github.com/fikrirnurhidayat/*
export GONOPROXY=github.com/fikrirnurhidayat/*
export GO111MODULE=on
export GOPRIVATE=github.com/fikrirnurhidayat/*

mod:
	go mod tidy
	go mod vendor

format:
	go fmt ./...

develop: format
	go run cmd/serve/main.go

build: format
	set -e mkdir -p target/bin
	go build -o target/bin/codeot-golang-executor ./cmd/serve/main.go

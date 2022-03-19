MAKEFLAGS += --silent

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

export GONOSUMDB=github.com/fikrirnurhidayat/*
export GONOPROXY=github.com/fikrirnurhidayat/*
export GO111MODULE=on
export GOPRIVATE=github.com/fikrirnurhidayat/*

setup:
	go mod vendor
	go install

develop:
	go run main.go

build:
	set -e mkdir out
	go build -o out/$(shell basename ${PWD}) main.go

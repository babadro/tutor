SHELL := /bin/bash

.PHONY: build run run-test-env run-test-env-with-bot test test-sleep down start-colima gen_mocks swag proto build-debug-binary

# check OS type
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
    GOARCH := amd64
else ifeq ($(UNAME_M),arm64)
    GOARCH := arm64
else ifeq ($(UNAME_M),aarch64)
	GOARCH := arm64
else
	GOARCH := unsupported_platform
endif


build:
	GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o release/app github.com/babadro/tutor/cmd/server

build-debug-binary:
	GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=0 go build -gcflags="all=-N -l" -o release/app github.com/babadro/tutor/cmd/server

run:
	docker-compose down -v && docker-compose build service && docker-compose up

# example: make run-test-env start-bot=true
run-test-env:
	docker-compose down -v && docker-compose build service && docker-compose up

run-test-env-with-bot:
	docker-compose down -v && docker-compose build service && docker-compose up

# example: make test filter=TestPolls
test:
	 (source .env.tests && go test ./... -p 1 -testify.m=$(filter) -v)

sleep-filter ?= TestPolls_Options
# example: make test-sleep sleep-filter=TestPolls_Options
test-sleep:
	(source .env.tests && SLEEP_MODE=true go test ./... -testify.m=$(sleep-filter) -v)

down:
	docker-compose down -v

start-colima:
	colima start -c 8 -m 8 --arch aarch64 --vm-type=vz --vz-rosetta --mount-type=virtiofs --vz-rosetta

gen_mocks:
	<placeholder>

swag:
	swagger generate server --exclude-main --server-package=internal/infra/restapi --model-package=internal/models/swagger -f swagger.yaml

dev-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
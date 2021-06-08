
LDFLAGS      := -w -s
MODULE       := github.com/figment-networks/kava-worker
VERSION_FILE ?= ./VERSION


# Git Status
GIT_SHA ?= $(shell git rev-parse --short HEAD)

ifneq (,$(wildcard $(VERSION_FILE)))
VERSION ?= $(shell head -n 1 $(VERSION_FILE))
else
VERSION ?= n/a
endif

all: build

.PHONY: plugin
plugin:
	CGO_ENABLED="1" go build -trimpath -o converter-plugin.so -buildmode=plugin ./cmd/converter-plugin

.PHONY: build
build: LDFLAGS += -X $(MODULE)/cmd/worker-kava/config.Timestamp=$(shell date +%s)
build: LDFLAGS += -X $(MODULE)/cmd/worker-kava/config.Version=$(VERSION)
build: LDFLAGS += -X $(MODULE)/cmd/worker-kava/config.GitSHA=$(GIT_SHA)
build:
      CGO_ENABLED=0 go build -o worker -ldflags '$(LDFLAGS)' ./cmd/worker-kava

.PHONY: pack-release
pack-release:
	@mkdir -p ./release
	@make build
	@mv ./worker ./release/worker
	@zip -r kava-worker ./release
	@rm -rf ./release

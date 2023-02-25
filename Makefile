SHELL=/usr/bin/env bash

all: build
.PHONY: all

unexport GOFLAGS

ldflags=-X=github.com/strahe/suialert/build.CurrentCommit=+git.$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

build: $(BUILD_DEPS)
	rm -f saas
	go build $(GOFLAGS) -o saas .
.PHONY: build

saas: build
.PHONY: saas
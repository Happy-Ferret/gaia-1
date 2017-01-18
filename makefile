UNAME := $(shell sh -c 'uname')
VERSION := $(shell sh -c 'git describe --always --tags')
ifdef GOBIN
PATH := $(GOBIN):$(PATH)
else
PATH := $(subst :,/bin:,$(GOPATH))/bin:$(PATH)
endif

# Standard Gaia build
default: prepare build

# Only run the build (no dependency grabbing)
build:
	go build -o gaia -ldflags \
		"-X main.Version=$(VERSION)" \
		./main.go

# Build with race detector
dev: prepare
	go build -race -o gaia -ldflags \
		"-X main.Version=$(VERSION)" \
		./main.go

# Build linux 64-bit, 32-bit and arm architectures
build-linux-bins: prepare
	GOARCH=amd64 GOOS=linux go build -o gaia_linux_amd64 \
								-ldflags "-X main.Version=$(VERSION)" \
								./main.go

# Get dependencies and use gdm to checkout changesets
prepare:
	go get ./...

sync:
	rsync -azvp --exclude '.git' --exclude '*.go' . p-axcoto:gaia

adhoc:
	ssh axcoto "cd gaia; ./run_linux.sh"

test-short:
	go test -short ./...

.PHONY: test-short

deploy: build push

GITHUB_USER=notyim
GITHUB_REPO=gaia
DESCRIPTION=$(shell sh -c 'git log --pretty=oneline | head -n 1')

UNAME := $(shell sh -c 'uname')
VERSION := $(shell sh -c 'git describe --always --tags')
CURRENT_VERSION := $(shell sh -c 'git rev-parse --short HEAD')

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

github-release:
	github-release release \
		--user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(CURRENT_VERSION) \
		--name "RELEASE $(CURRENT_VERSION)" \
		--description "$(DESCRIPTION)"

	github-release upload \
		--user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(CURRENT_VERSION) \
		--name "gaia-linux" \
		--file gaia_linux_amd64

clean-influx:
	echo "Clean influxdb"

release: build-linux-bins github-release

# Production task
ssh-deploy:
	ssh noty "curl https://github.com/NotyIm/gaia/releases/download/$(CURRENT_VERSION)/gaia-linux -o /var/app/gaia/bin/gaia; systemctl restart gaia"

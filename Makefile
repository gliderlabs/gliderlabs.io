## PHONY

.PHONY: help dev setup build ui image clean clobber testenv deploy linux darwin js css

## VARIABLES

BIN := glio

## COMMANDS

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: deps-update css ## run development runner
	comlab dev

setup: deps-update ## setup development environment

build: linux darwin ## build linux and darwin binaries

ui: css ## build static css and js artifacts

image: static-bin ## build docker container
	docker build -t gliderlabs/glio .

clean: ## delete typical build artifacts
	rm -rf build

clobber: clean ## reset dev environment
	rm -rf vendor
	rm -rf ui/node_modules
	rm -rf ui/static/semantic
	rm -f .git/deps-*

## TESTS

test: test-go ## run common tests

test-all: test-go test-env ## run ALL tests

test-go: ## test golang packages
	go test -race ./pkg/... ./cmd/... ./com/...

test-env: ## test dev environment
	docker build -t glio-env -f dev/setup/Dockerfile .
	docker rmi glio-env

## DEPENDENCIES

deps-update: ## update dependencies if changed
	./dev/deps.sh

deps-go:
	glide install
	git log -n 1 --pretty=format:%h -- glide.yaml > .git/deps-go

deps-css:
	make -C ui/semantic
	git log -n 1 --pretty=format:%h -- ui/semantic > .git/deps-css


## DEPLOY

deploy: clean static-bin ## deploy to gliderlabs.io
	heroku container:push web -a gliderlabs-io

## ALIASES

linux: build/linux/${BIN} ## build linux binary

darwin: build/darwin/${BIN} ## build darwin binary

static-bin: build/linux-static/${BIN} ## build static linux binary

css: ui/static/semantic ## build css from semantic ui

## PATHS

build/linux-static/${BIN}:
	mkdir -p build/linux-static
	GOOS=linux CGO_ENABLED=0 go build -a \
		-installsuffix cgo \
		-o build/linux-static/${BIN} ./cmd/${BIN}

build/linux/${BIN}:
	mkdir -p build/linux
	GOOS=linux go build -o build/linux/${BIN} ./cmd/${BIN}

build/darwin/${BIN}:
	mkdir -p build/darwin
	GOOS=darwin go build -o build/darwin/${BIN} ./cmd/${BIN}

ui/static/semantic: ui/semantic
	make deps-css

ui/semantic:
	mv ui/semantic.ignore ui/semantic

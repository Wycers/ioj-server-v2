.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v0.0.2

APP 			= ioj-server
SERVER_BIN  	= ./dist/${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= 1# $(shell git rev-list --all --count)
GIT_HASH        = asd#$(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

build:
	go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./cmd/${APP}

start:
	go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./cmd/${APP} -f ./configs/server.yaml

swagger:
	swag init --parseDependency --generalInfo ./cmd/${APP}/main.go --output ./internal/app/swagger

wire:
	wire gen ./cmd/${APP}

air-build: wire swagger build

test:
	@go test -v $(shell go list ./...)

clean:
	rm -rf data release $(SERVER_BIN) internal/app/test/data cmd/${APP}/data

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}

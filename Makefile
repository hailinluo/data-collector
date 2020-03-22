GO = go
GOGET = $(GO) get
BUILD = ${GO} build ${GOFLAGS}
BIN = data-collector
GOBUILD = ${BUILD} -o ./bin/${BIN} ./main.go

all: build

.PHONY: build
build:
	mkdir -p bin
	${GOBUILD}

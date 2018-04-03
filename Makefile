GO=go
GOVENDOR=govendor
SRCROOT=$(shell pwd)
RELEASE=$(shell git describe --abbrev=0 --tags)
COMMIT=$(shell git rev-parse --short HEAD)

SOURCES=config.go error.go ping.go server.go http.go health.go

.PHONY: all container clean vendor

all: kube-health

clean:
	rm -f kube-health

container: kube-health
	docker build -t kube-health ${SRCROOT}

vendor:
	cd src/health; ${GOVENDOR} update +v +m

kube-health: ${SOURCES}
	${GO} build -tags netgo -ldflags "-X main.version=${RELEASE} -X main.revision=${COMMIT}" -o $@ $^

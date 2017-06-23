GO=go
GOVENDOR=govendor
GOPATH=$(shell pwd)

SOURCES=src/health/config.go src/health/error.go src/health/ping.go src/health/server.go src/health/http.go

.PHONY: all container clean vendor

all: bin/kubehealth

clean:
	rm -f bin/kubehealth

container: bin/kubehealth
	docker build -t kube-health ${GOPATH}

vendor:
	cd src/health; ${GOVENDOR} update +v +m

bin/kubehealth: src/health.go ${SOURCES}
	${GO} build -o $@ $<

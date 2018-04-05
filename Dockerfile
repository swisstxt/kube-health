# Create a Docker container that reports health from inside a Kubernetes cluster


# Build stage
FROM golang:1.10-alpine AS build

# Pass build arguments
ARG RELEASE=
ARG COMMIT=

# Copy source code
WORKDIR /go/src/kube-health/
COPY . .

# Build container
RUN CGO_ENABLED=0 go build -tags netgo -ldflags "-X main.version=$RELEASE -X main.revision=$COMMIT" -o kube-health *.go


# Image creation stage
FROM alpine:latest

# Some metadata
LABEL name="kube-health"
LABEL version="1.0"
LABEL maintainer "gregor.riepl@swisstxt.ch"
LABEL description="Simple Docker container that runs a configurable set of \
health checks inside a Kubernetes pod and exposes the result over HTTP."

ENV LANG C.UTF-8

RUN apk update && apk upgrade

# Install dependencies
RUN apk --no-cache add ca-certificates iputils

RUN adduser -h /health -H -D -s /bin/bash health && \
	mkdir -p /health && \
	chown health /health && \
	chmod 0750 /health

# Install the web server
COPY --from=build /go/src/kube-health/kube-health /bin/kubehealth
# Install the default configuration file
COPY example-config.json /etc/kubehealth/config.json

# Can run with reduced privileges
USER health

# This is the port our web server will listen on
EXPOSE 8080
# The main entry point for this container
CMD ["/bin/kubehealth", "/etc/kubehealth/config.json"]

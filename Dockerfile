# Create a Docker container that reports health from inside a Kubernetes cluster

FROM debian:stable-slim

# Some metadata
LABEL name="kube-health"
LABEL version="1.0"
LABEL maintainer "gregor.riepl@swisstxt.ch"
LABEL description="Simple Docker container that runs a configurable set of \
health checks inside a Kubernetes pod and exposes the result over HTTP."

ENV LANG C.UTF-8

RUN apt-get update && apt-get upgrade -y

# Uncomment if you need more locales than the default C.UTF-8
# RUN apt-get install -y locales && \
# 	rm -rf /var/lib/apt/lists/* && \
# 	localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
# ENV LANG en_US.UTF-8
# RUN apt-get update

# Create a user for the web server
RUN useradd -d /health -r -M -s /bin/bash health && \
	mkdir -p /health && \
	chown health /health && \
	chmod 0750 /health

# Install the web server
COPY bin/kubehealth /usr/bin/kubehealth
# Install the default configuration file
COPY example-config.json /etc/kubehealth/config.json

# TODO listening for ICMP packets requires root
#USER health

# This is the port our web server will listen on
EXPOSE 8080
# The main entry point for this container
CMD ["/usr/bin/kubehealth", "/etc/kubehealth/config.json"]

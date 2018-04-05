# kube-health

Simple, configurable web server for running health checks.

Can be built into a Docker container - very useful for running
health checks in a Kubernetes cluster.

## Build

Get [Git](https://git-scm.com/) and [Go](https://golang.org/).

Fetch and build:
```
go get github.com/swisstxt/kube-health
```

## Use

Copy the [example configuration](example-config.json) to `config.json`.

Modify the tests in `config.json` to your heart's content.

The log levels are as follows:

0. only return status and no check results
1. return check results when critical
2. also return check results when unknown or warning
3. always return check results

It is recommended to keep the log level at 1 or 2.

Save and run:
```
kube-health config.json
```

Open http://127.0.0.1:8080/ in your web browser and see the result of the example checks.

## Install

Get [Docker](https://www.docker.com/).

Build the container:
```
cd $GOPATH/src/github.com/swisstxt/kube-health
docker build -t kube-health --build-arg "RELEASE=$(git describe --abbrev=0 --tags --exact-match)" --build-arg "COMMIT=$(git rev-parse --short HEAD)"
```

Test the container:
```
docker run -ti --rm -p 127.0.0.1:8080:8080 kube-health
```

Connect to http://localhost:8080 as above.

Tag and push the docker container to your Docker registry (replace docker-registry.tld with the
correct host name):
```
docker tag kube-health <docker-registry.tld>/kube-health
docker push <docker-registry.tld>/kube-health
```

You should also tag the images with the release version, if you are build a release.

## Kubernetes

Copy the [example configuration](example-config.json) to `config.json`.

Edit kube-health.yaml and change the line that says `- image: kube-health` to say:
`- image: <docker-registry.tld>/kube-health`.

Create a deployment in your Kubernetes cluster and add the configuration file:
```
kubectl apply -f kube-health.yaml -n kube-system
kubectl create configmap kube-health --from-file config.json -n kube-system
```

## Legal

kube-health is copyright © 2017-2018 by Gregor Riepl and Swiss TXT AG.

You may use it under the terms of the MIT license.
Please see the [LICENSE](LICENSE) file for details.

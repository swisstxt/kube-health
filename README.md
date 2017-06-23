# kube-health

Simple, configurable web server for running health checks.

Can be built into a Docker container - very useful for running
health checks in a Kubernetes cluster.

## Use

Get [Git](https://git-scm.com/), [Go](https://golang.org/),
[GNU make](https://www.gnu.org/software/make/) and [Docker](https://www.docker.com/).

Get the sources:
```
git clone https://github.com/swisstxt/kube-health
```

Build the container:
```
make
```

Run the container:
```
docker run -ti --rm kube-health
```

## Install

Push the docker container to your Docker registry:
```
TODO
```

Copy the [example configuration](example-config.json) to `config.json`

Create a deployment in your Kubernetes cluster and add the configuration file:
```
kubectl apply -f kube-health.yaml -n default
kubectl create configmap kube-health --from-file config.json
```

## Customise

Open `config.json` in your favourite text editor.

Modify the tests to your heart's content.

Save and run:
```
kubectl create configmap kube-health --from-file config.json -o yaml --dry-run | kubectl apply -f - -n default
```

## Legal

kube-health is copyright © 2017 by Gregor Riepl and Swiss TXT AG.

You may use it under the terms of the MIT license.
Please see the [LICENSE](LICENSE) file for details.

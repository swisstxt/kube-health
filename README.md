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
docker run -ti --rm -p 127.0.0.1:8080:8080 kube-health
```

Open http://127.0.0.1:8080/ in your web browser and see the result of the example checks.

## Install

Tag and push the docker container to your Docker registry (replace docker-registry.tld with the
correct host name):
```
docker tag kube-health <docker-registry.tld>/kube-health
docker push <docker-registry.tld>/kube-health
```

Copy the [example configuration](example-config.json) to `config.json`.

Edit kube-health.yaml and change the line that says `- image: kube-health` to
`- image: <docker-registry.tld>/kube-health`.

Create a deployment in your Kubernetes cluster and add the configuration file:
```
kubectl apply -f kube-health.yaml -n kube-system
kubectl create configmap kube-health --from-file config.json -n kube-system
```

## Customise

Open `config.json` in your favourite text editor.

Modify the tests to your heart's content.

Save and run:
```
kubectl create configmap kube-health --from-file config.json -o yaml --dry-run | kubectl apply -f - -n kube-system
```

## Legal

kube-health is copyright © 2017 by Gregor Riepl and Swiss TXT AG.

You may use it under the terms of the MIT license.
Please see the [LICENSE](LICENSE) file for details.

# Installation Guide

### Prerequisite

* redis installed:

```text
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis --version 20.7.0
```

### Build and Install event-consumer

#### Build docker image:

```text
cd event-consumer/
```

```text
docker build $PWD --file docker/Dockerfile --tag event-consumer-image:0.0.1
```

```text
minikube image load event-consumer-image:0.0.1
```

#### Deploy helm chart:

```text
helm install event-consumer chart/
```

#### Uninstall:

```text 
helm uninstall event-consumer
```

```text
minikube image rm docker.io/library/event-consumer-image:0.0.1
```

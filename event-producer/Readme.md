# Installation Guide

### Prerequisite

* redis installed:

```text
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis --version 20.7.0
```

### Build and Install event-producer

### Build docker image:

```text
cd event-producer/
```

```text
docker build $PWD --file docker/Dockerfile --tag event-producer-image:0.0.1
```

```text
minikube image load event-producer-image:0.0.1
```

#### Deploy helm chart:

```text
helm install event-producer chart/
```

### Uninstall:

```text 
helm uninstall event-producer
```

```text
minikube image rm docker.io/library/event-producer-image:0.0.1
```
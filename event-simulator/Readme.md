# Installation Guide

### Build docker image:

```text
cd event-simulator/
```

```text
docker build $PWD --file docker/Dockerfile --tag event-simulator-image:0.0.1
```

```text
minikube image load event-simulator-image:0.0.1
```

### Deploy helm chart:

```text
helm install event-simulator chart/
```

### Uninstall:

```text 
helm uninstall event-simulator
```

```text
minikube image rm docker.io/library/event-simulator-image:0.0.1
```
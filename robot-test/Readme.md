# Installation Guide

### Build docker image:

```text
cd robot-test/
```

```text
docker build $PWD --file docker/Dockerfile --tag robot-test-image:0.0.1
```

```text
minikube image load robot-test-image:0.0.1
```

### Deploy helm chart:

```text
helm install robot-test chart/
```

### Manually execute test case inside robot-test pod:

Execute functional tests:

```text
robot test_functional.robot
```

Execute robustness tests:

```text
robot test_robustness.robot
```

## Uninstall:

```text 
helm uninstall robot-test
```

```text
minikube image rm docker.io/library/robot-test-image:0.0.1
```
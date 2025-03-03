# Installation Guide

This document shows the steps required to deploy this project on the kubernetes cluster.

### Prerequisites

* Docker running so you can build the docker images locally.
* Kubernetes Cluster running for example: Minikube. 
  * Docker images will be loaded into Minikube.
* Helm installed so you can build, download dependencies and install the helm charts.

## 1. Building the docker images

Following images must be created and loaded to the minikube cluster:

- event-producer-image
- event-consumer-image
- event-simulator-image

#### 1. Build docker images:

```text
cd event-producer/
docker build $PWD --file docker/Dockerfile --tag event-producer-image:0.0.1
```

```text
cd event-consumer/
docker build $PWD --file docker/Dockerfile --tag event-consumer-image:0.0.1
```

```text
cd event-simulator/
docker build $PWD --file docker/Dockerfile --tag event-simulator-image:0.0.1
```

#### 2. Load the images to minikube cluster:

```text
minikube image load event-producer-image:0.0.1
```

```text
minikube image load event-consumer-image:0.0.1
```

```text
minikube image load event-simulator-image:0.0.1
```

#### 3. Verify images have been loaded into minikube

```text
minikube image list
```

* How to remove an image from minikube if required:

```text
minikube image rm docker.io/library/event-producer-image:0.0.1
```

## 2. Build the helm chart

The `chart/charts/` folder contains the helm charts that will be installed. 

It contains 4 charts in total, 3 internal charts i.e `event-consumer-0.0.1.tgz`, `event-producer-0.0.1.tgz` and `event-simulator-0.0.1.tgz`
and 1 external chart that will be downloaded from an artifacotry i.e `redis-20.7.0.tgz`. 

Two microservices the `event-producer` and the `event-consumer` have a dependency on the `redis` microservice in this project.
To make the installation process simpler the redis microservice has been added to the dependencies in the `chart/Chart.yaml` file. 
This means that you don't have to manually install redis microservice as it will be packaged alongside the other microservices 
  i.e. `event-producer`, `event-consumer`, `event-simulator`. 

- Link to the dependency used: https://artifacthub.io/packages/helm/bitnami/redis

The following redis dependency is defined in `chart/Chart.yaml`:

```yaml
dependencies:
  - name: redis
    version: 20.7.0
    repository: https://charts.bitnami.com/bitnami
  - name: event-producer
    version: 0.0.1
    repository: "file://../event-producer/chart"
    condition: event-producer.enabled
  - name: event-consumer
    version: 0.0.1
    repository: "file://../event-consumer/chart"
    condition: event-consumer.enabled
  - name: event-simulator
    version: 0.0.1
    repository: "file://../event-simulator/chart"
    condition: event-simulator.enabled
```

#### 1. Run the following helm command to download the redis dependency:

```text
helm dependency update chart/
```

- Helm chart will be downloaded under: `chart/charts/redis-20.7.0.tgz`.


#### 2. Build the helm chart archive

```text
helm package chart/
```

- Helm chart archive will be created: `application-1.0.0.tgz`. It will contain all 4 helm charts,
i.e. `event-producer-0.0.1.tgz`, `event-consumer-0.0.1.tgz`, `event-simulator-0.0.1.tgz`, `redis-20.7.0.tgz`.


## 3. Installation

#### 1. Install the helm chart:

```text
helm install application application-1.0.0.tgz
```

The following kubernetes resources will be created:

```text
kubectl get all
```

| **Pods**                        | additional information                                                                                                 |
|---------------------------------|------------------------------------------------------------------------------------------------------------------------|
| event-consumer-7c9f4b48f5-wv5kg | consumes a Notification Message (NM) from redis stream (NotificationStream).                                           |
| event-producer-74b59dd7d-6lt8n  | produces a Notification Messages (NM) to redis stream (NotificationStream) after a Event Notification (EM) is created. |
| event-simulator-f9db6d494-7sdp7 | posts an Event Message (EM) to event-producer.                                                                         |
| redis-master-0                  | the `event-consumer` and `event-producer` has a dependency on this pod. <br/>Pod contains`redis-cli`.                  |
| redis-replicas-0                | -                                                                                                                      |
| redis-replicas-1                | -                                                                                                                      |
| redis-replicas-2                | -                                                                                                                      |

| **Services**       | additional information                                                       |
|--------------------|------------------------------------------------------------------------------|
| event-producer-svc | the `event-consumer` has a dependency on this `service`                      |
| redis-master       | the `event-consumer` and `event-producer` has a dependency on this `service` |

| **Secrets** | additional information                                                           |
|-------------|----------------------------------------------------------------------------------|
| redis       | the `event-consumer` and `event-producer` has a helm dependency on this `secret` |

#### 2. Verify Redis stream has been created

#### 1. Get the redis secret:

```text
kubectl get secrets redis -o jsonpath='{.data.redis-password}' | base64 -d
```

* Example output - value will be different:

```text
beLzE8UmFN
```

#### 2. Exec into redis pod:

```text
kubectl exec -it redis-master-0 -- /bin/bash
```

#### 3. Log into redis using redis-cli and use the password in step 1.:
- https://redis.io/docs/latest/operate/rs/references/cli-utilities/redis-cli/

```text
redis-cli -p 6379 -a <secret>
```

#### 4. redis-cli command to get list of keys of existing streams:
- https://redis.io/docs/latest/commands/keys/

```text
keys *
```

* Output:

```text
1) "NotificationStream" 
```

##### 5. redis-cli command to return information about the stream:
- https://redis.io/docs/latest/commands/xinfo-stream/

```text
XINFO STREAM NotificationStream
```

* Output:

```text
 1) "length"
 2) (integer) 100
 3) "radix-tree-keys"
 4) (integer) 4
 5) "radix-tree-nodes"
 6) (integer) 11
 7) "last-generated-id"
 8) "1740523293390-0"
 9) "max-deleted-entry-id"
10) "0-0"
11) "entries-added"
12) (integer) 300990
13) "recorded-first-entry-id"
14) "1740522183181-0"
15) "groups"
16) (integer) 1
17) "first-entry"
18) 1) "1740522183181-0"
    2) 1) "EventStream"
       2) "{\"id\":3082736038337373727,\"description\":\"hello world\",\"timestamp\":\"2025-02-25T22:23:03.181361595Z\"}"
19) "last-entry"
20) 1) "1740523293390-0"
    2) 1) "CreateStream"
       2) "stream created"
```

## 5. Uninstall the helm chart:

```text
helm uninstall application
```

Remove the docker images from minikube:

```text
minikube image rm docker.io/library/event-producer-image:0.0.1
```

```text
minikube image rm docker.io/library/event-consumer-image:0.0.1
```

```text
minikube image rm docker.io/library/event-simulator-image:0.0.1
```

import time
from typing import Optional, Any

from kubernetes import client, config
from kubernetes.stream import stream

import logging

logger = logging.getLogger(__name__)

def _load_kube_config():
    is_in_cluster = True
    if is_in_cluster:
        config.load_incluster_config()
    else:
        config.load_kube_config("~/.kube/config")

def get_list_of_pods(namespace: str):
    _load_kube_config()
    api_instance = client.CoreV1Api()

    resp = api_instance.list_namespaced_pod(namespace)
    for i in resp.items:
        print(f"{i.status.pod_ip}\t{i.metadata.namespace}\t{i.metadata.name}")
    return resp

def get_pod_name_on_label(namespace: str, label: Optional[str] = None) -> list[str]:
    _load_kube_config()
    api_instance = client.CoreV1Api()

    pods = api_instance.list_namespaced_pod(namespace, label_selector=label)
    name = [str(pod.metadata.name) for pod in pods.items]
    return name

def get_secret_data(namespace: str, secret_name: str) -> str:
    _load_kube_config()
    api_instance = client.CoreV1Api()

    secret = api_instance.read_namespaced_secret(name=secret_name, namespace=namespace)
    return secret.data

def pod_exec_command(namespace: str, exec_command: [], pod_name: str) -> Any:
    _load_kube_config()
    api_instance = client.CoreV1Api()

    resp = stream(api_instance.connect_get_namespaced_pod_exec,
                  name=pod_name,
                  namespace=namespace,
                  command=exec_command,
                  stderr=True, stdin=False,stdout=True, tty=False)
    return resp

def pod_delete(namespace: str, name: str):
    _load_kube_config()
    api_instance = client.CoreV1Api()
    resp = api_instance.delete_namespaced_pod(name=name, namespace=namespace)
    #print("Deleted pod: ",resp)


def test_print(namespace: str, name: str):
    _load_kube_config()
    api_instance = client.CoreV1Api()
    pod = api_instance.read_namespaced_pod(name=name, namespace=namespace)


    print(f"TEST PRINT hello: {pod.status.phase}")


def verify_pod_is_down(namespace: str, pod_name: str) -> bool:
    _load_kube_config()
    api_instance = client.CoreV1Api()
    attempt = 1
    while attempt <= 10:
        pod = api_instance.read_namespaced_pod(name=pod_name, namespace=namespace)
        pod_state = pod.status.phase
        pod_ready = pod.status.container_statuses[0].ready
        print(f"Verify pod down attempt: {attempt}, pod state: {pod_state}, pod ready: {pod_ready}")

        if not pod_ready:
            return True

        attempt += 1
        time.sleep(1)
    raise Exception("Pod failed to restart")

def verify_pod_is_up(namespace: str, name: str) -> bool:
    api_instance = client.CoreV1Api()
    attempt = 1
    while attempt <= 10:
        pod = api_instance.read_namespaced_pod(name=name, namespace=namespace)
        pod_state = pod.status.phase
        pod_ready = pod.status.container_statuses[0].ready

        print(f"Verify pod up attempt: {attempt}, pod state: {pod_state}, pod ready: {pod_ready}")

        if (pod_state == 'Running') & pod_ready:
           return True

        attempt += 1
        time.sleep(5)
    raise Exception("Pod failed to start")
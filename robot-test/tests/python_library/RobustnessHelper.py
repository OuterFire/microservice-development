from helper import kubernetes_helper

def wait_for_pod_restart(namespace: str, pod: str) -> bool:
    pod_restarted = kubernetes_helper.verify_pod_is_down(namespace, pod)
    print("pod restarted",pod_restarted)

    status = kubernetes_helper.verify_pod_is_up(namespace=namespace, name=pod)
    print("pod restarted",status)
    return status


def restart_pod(namespace: str, pod_label: str) -> str:
    pod = kubernetes_helper.get_pod_name_on_label(namespace=namespace, label=pod_label)
    if len(pod) != 1:
        raise Exception(f"Only 1 pod should exist: {pod}")

    kubernetes_helper.pod_delete(namespace=namespace, name=pod[0])
    return pod[0]

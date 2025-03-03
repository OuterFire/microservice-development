from typing import List
from helper import kubernetes_helper

import ast
import base64
import json

def get_notification_message_from_redis(namespace: str, redis_pod_label: str, secret_name: str, stream_name: str) -> List[str]:
    redis_pod = kubernetes_helper.get_pod_name_on_label(namespace=namespace, label=redis_pod_label)
    if len(redis_pod) != 1:
        raise Exception(f"Only 1 redis pod should exist in {redis_pod}")

    redis_secret = get_redis_secret(namespace, secret_name)

    grep_cmd = 'grep {'
    tail_cmd = 'tail -n1'
    redis_cli_cmd = f"redis-cli --no-auth-warning --pass {redis_secret} -p 6379 XRANGE {stream_name} - + | {grep_cmd} | {tail_cmd}"
    exec_command = ["bash", "-c", redis_cli_cmd]

    print("Exec redis cli command:", exec_command)

    resp = kubernetes_helper.pod_exec_command(namespace=namespace, exec_command=exec_command, pod_name=redis_pod[0])
    notification_message_list = resp.splitlines()

    if len(notification_message_list) == 0:
        raise Exception(f"Got no notification message from redis, response: {notification_message_list}")

    notification_messages = []
    for i in range(len(notification_message_list)):
        print(notification_message_list[i])
        msg = ast.literal_eval(notification_message_list[i])
        notification_messages.append(msg)
    return notification_messages

def demo(namespace: str, redis_pod_label: str):
    redis_pod = kubernetes_helper.get_pod_name_on_label(namespace=namespace, label=redis_pod_label)
    if len(redis_pod) != 1:
        raise Exception(f"Only 1 redis pod should exist in {redis_pod}")
    kubernetes_helper.test_print(namespace=namespace, name=redis_pod[0])


def get_redis_secret(namespace: str, secret_name: str) -> str:
    secret = kubernetes_helper.get_secret_data(namespace=namespace, secret_name=secret_name)
    x = json.loads(json.dumps(secret))
    return base64.b64decode(x["redis-password"]).decode("utf-8")

from __future__ import print_function

import logging
import random
import time

import grpc
from proto import todolist_pb2
from proto import todolist_pb2_grpc
from flask import Flask

app = Flask(__name__)


def run():

    port = "5000"

    with grpc.insecure_channel('172.31.33.200:' + str(port)) as channel:
        stub = todolist_pb2_grpc.TodoListStub(channel)
        metadata = (('dapr-app-id', 'server'),)
        response = stub.GetTodolist(request=todolist_pb2.google_dot_protobuf_dot_empty__pb2.Empty(), metadata=metadata)
        return response


@app.route("/todo")
def hello_world():
    return str(run())


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=3000)

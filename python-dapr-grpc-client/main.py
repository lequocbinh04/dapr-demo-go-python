from __future__ import print_function

import logging
import random
import time

import grpc
from proto import todolist_pb2
from proto import todolist_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:3030') as channel:
        stub = todolist_pb2_grpc.TodoListStub(channel)
        metadata = (('dapr-app-id', 'server'),)
        response = stub.GetTodolist(request=todolist_pb2.google_dot_protobuf_dot_empty__pb2.Empty(), metadata=metadata)
        print(response)


if __name__ == '__main__':
    run()
import logging
import time
import queue

import grpc
from locust import FastHttpUser, HttpUser, task

from user import user_pb2_grpc

from utils import (
    generate_user,
    generate_user_address,
    generate_user_address_request,
    generate_user_request,
    grpc_stopwatch,
    GrpcUser,
    LocustHttpUser,
)

from gevent import monkey
import grpc.experimental.gevent as grpc_gevent

# patch grpc so that it uses gevent instead of asyncio
monkey.patch_all()
grpc_gevent.init_gevent()

logger = logging.getLogger(__name__)

BATCH_SIZE = 10


class GRPCLoadTest(GrpcUser):
    stub_class = user_pb2_grpc.UserServiceStub

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUser")
    def create_user_unary(self):
        return self.stub.CreateUser(generate_user_request(generate_user()))

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserAddress")
    def create_user_address_unary(self):
        return self.stub.CreateUserAddress(
            generate_user_address_request(generate_user_address())
        )


class GRPCStreamLoadTest(GrpcUser):
    compression = grpc.Compression.Gzip
    stub_class = user_pb2_grpc.UserServiceStub

    user_queue = queue.Queue()
    user_address_queue = queue.Queue()
    sender_user = None
    sender_user_address = None

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.sender_user = self.stub.CreateUserBidirectional(
            iter(self.user_queue.get, None),
        )
        self.sender_user_address = self.stub.CreateUserAddressBidirectional(
            iter(self.user_address_queue.get, None),
        )

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserBidirectional")
    def create_user_bidirectional_stream(self):
        self.user_queue.put(generate_user_request())
        return next(self.sender_user)

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserAddressBidirectional")
    def create_user_address_bidirectional_stream(self):
        self.user_address_queue.put(generate_user_address_request())
        return next(self.sender_user_address)


class FastHttpLoadTest(LocustHttpUser, FastHttpUser):

    @task
    def create_user(self):
        self.client.post("/users", json=generate_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=generate_user_address())


class HttpLoadTest(LocustHttpUser, HttpUser):

    @task
    def create_user(self):
        self.client.post("/users", json=generate_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=generate_user_address())

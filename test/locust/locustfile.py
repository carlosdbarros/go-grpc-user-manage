# test/locust/locustfile.py
import logging
import time
import queue
import threading
from functools import wraps

import grpc
from locust import FastHttpUser, HttpUser, User, task, events, between
from locust.exception import LocustError

from user import user_pb2_grpc, user_pb2

from utils import (
    GrpcUser, generate_user,
    generate_user_address,
    generate_user_address_request,
    generate_user_request,
    generate_user_iterator,
)

from gevent import monkey
import grpc.experimental.gevent as grpc_gevent

# patch grpc so that it uses gevent instead of asyncio
monkey.patch_all()
grpc_gevent.init_gevent()

logger = logging.getLogger(__name__)

BATCH_SIZE = 10


def grpc_stopwatch(rpc_name=None):
    def decorator(fn):
        @wraps(fn)
        def wrapper(self, *args, **kwargs):
            if self.channel_closed:
                return
            response = None
            exception = None
            start_time = time.time()
            start_perf_counter = time.perf_counter()

            try:
                response = fn(self, *args, **kwargs)
            except grpc.RpcError as e:
                exception = e

            response_time = (time.perf_counter() - start_perf_counter) * 1000
            response_size = response.ByteSize() if response else 0
            name = rpc_name.split("/")[-1] if rpc_name else fn.__name__

            self.environment.events.request.fire(
                request_type="gRPC",
                response_time=response_time,
                name=name,
                response=response,
                response_length=response_size,
                start_time=start_time,
                url=rpc_name,
                exception=exception,
                context={},
            )
            return response

        return wrapper

    return decorator


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


class GrpcStreamLoadTest(GrpcUser):
    stub_class = user_pb2_grpc.UserServiceStub

    _lock = threading.Lock()
    sender_queue = queue.Queue()
    sender = None

    def on_start(self):
        super().on_start()
        self.sender = self.stub.CreateUserStreamStream(iter(self.sender_queue.get, None))

    def on_stop(self, force=False):
        self.sender_queue.put(None)
        time.sleep(1)
        self.sender.cancel()
        self.sender = None
        super().on_stop(force=True)

    @task
    def create_user_stream_stream(self):
        if self.channel_closed:
            return
        response = None
        exception = None
        start_time = time.time()
        start_perf_counter = time.perf_counter()

        # with self._lock:
        try:
            request = generate_user_request()
            self.sender_queue.put(request)
            response = next(self.sender)
        except grpc.RpcError as e:
            exception = e

        response_time = (time.perf_counter() - start_perf_counter) * 1000
        response_size = response.ByteSize() if response else 0

        self.environment.events.request.fire(
            request_type="gRPC",
            response_time=response_time,
            name="CreateUserStreamStream",
            response=response,
            response_length=response_size,
            start_time=start_time,
            url="grpc.user.UserService/CreateUserStreamStream",
            exception=exception,
            context={},
        )


class FastHttpLoadTest(FastHttpUser):
    # min_wait = 1000
    # max_wait = 10000
    host = "http://localhost:8080"

    def on_stop(self, force=False):
        super().stop(force=True)

    @task
    def create_user(self):
        self.client.post("/users", json=generate_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=generate_user_address())

    # @task
    # def create_user_batch(self):
    #     for _ in range(BATCH_SIZE):
    #         self.client.post("/users", json=generate_user())


class HttpLoadTest(HttpUser):
    # min_wait = 1000
    # max_wait = 10000
    host = "http://localhost:8080"

    @task
    def create_user(self):
        self.client.post("/users", json=generate_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=generate_user_address())

# test/locust/locust_grpc.py
import logging
import time
from functools import wraps

import grpc
from locust import FastHttpUser, HttpUser, User, task, events, between
from locust.exception import LocustError

from user import user_pb2_grpc, user_pb2

from utils import (
    generate_user,
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


class GrpcUser(User):
    abstract = True
    min_wait = 1000
    max_wait = 10000

    host = None
    stub_class = None

    def __init__(self, environment):
        super().__init__(environment)
        for attr in ["host", "stub_class"]:
            if not getattr(self, attr):
                raise LocustError(
                    f"Class {self.__class__.__name__} missing required attribute {attr}"
                )

        self._channel = grpc.insecure_channel(self.host.lstrip("http://"))
        # grpc.channel_ready_future(self._channel).result(timeout=15)
        self.stub = self.stub_class(self._channel)
        # interceptor = LocustClientInterceptor(environment=environment)
        # self._channel = grpc.intercept_channel(self._channel, interceptor)
        # self.stub = self.stub_class(self._channel)
        self._channel_closed = False

    def on_start(self):
        self._channel_closed = False

    def on_stop(self, force=False):
        self._channel_closed = True
        time.sleep(1)
        self._channel.close()
        super().stop(force=True)


def grpc_decorator(rpc_name=None):
    def decorator(fn):
        @wraps(fn)
        def wrapper(self, *args, **kwargs):
            if self._channel_closed:
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
    host = "localhost:50051"
    stub_class = user_pb2_grpc.UserServiceStub

    @task
    @grpc_decorator("grpc.user.UserService/CreateUser")
    def create_user_unary(self):
        return self.stub.CreateUser(generate_user_request(generate_user()))

    @task
    @grpc_decorator("grpc.user.UserService/CreateUserAddress")
    def create_user_address_unary(self):
        return self.stub.CreateUserAddress(
            generate_user_address_request(generate_user_address())
        )


class HttpLoadTest(HttpUser):
    min_wait = 1000
    max_wait = 10000
    host = "http://localhost:8080"

    # @task
    # def create_user(self):
    #     self.client.post("/users", json=generate_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=generate_user_address())


class FastHttpLoadTest(FastHttpUser):
    min_wait = 1000
    max_wait = 10000
    host = "http://localhost:8080"

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

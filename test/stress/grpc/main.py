import logging
import time
import string
import grpc

from typing import Any, Callable

import grpc.experimental.gevent as grpc_gevent
from grpc_interceptor import ClientInterceptor

from locust import User, task
from locust.exception import LocustError

from user import user_pb2, user_pb2_grpc


logger = logging.getLogger(__name__)


class UserServiceServicer(user_pb2_grpc.UserServiceServicer):
    def CreateUser(self, request, context):
        user = user_pb2.User(
            id="1",
            name=request.name,
            email=request.email,
        )
        return user


# patch grpc so that it uses gevent instead of asyncio
grpc_gevent.init_gevent()


class LocustInterceptor(ClientInterceptor):
    def __init__(self, environment, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.env = environment

    def intercept(self, method: Callable, request_or_iterator: Any, call_details: grpc.ClientCallDetails):
        response = None
        exception = None
        start_perf_counter = time.perf_counter()
        response_length = 0
        try:
            response = method(request_or_iterator, call_details)
            response_length = response.result().ByteSize()
        except grpc.RpcError as e:
            exception = e
        else:
            response_time = (time.perf_counter() - start_perf_counter) * 1000
            self.env.events.request.fire(
                request_type="grpc",
                name=call_details.method,
                response_time=response_time,
                response_length=response_length,
                response=response,
                context=None,
                exception=exception,
            )
        return response


class gRPCClient(User):
    abstract = True
    stub_class = None

    def __init__(self, environment):
        super().__init__(environment)
        required_attrs = ["host", "stub_class"]
        missing_attrs = [attr for attr in required_attrs if not hasattr(self, attr)]
        if missing_attrs:
            raise LocustError(f"Required attributes {missing_attrs} not set on {self.__class__.__name__}")

    def on_start(self):
        super().on_start()
        self._channel = grpc.insecure_channel(self.host)
        interceptor = LocustInterceptor(environment=self.environment)
        self._channel = grpc.intercept_channel(self._channel, interceptor)
        self.stub = user_pb2_grpc.UserServiceStub(self._channel)

    def on_stop(self):
        self._channel.close()
        super().on_stop()


def generate_random_string(length):
    return "".join([string.ascii_letters[i % len(string.ascii_letters)] for i in range(length)])


class gRPCUserClient(gRPCClient):
    host = "localhost:50051"
    stub_class = user_pb2_grpc.UserServiceStub

    @task
    def CreateUser(self):
        request = user_pb2.CreateUserRequest(
            name=generate_random_string(10),
            email=generate_random_string(10) + "@gmail.com",
            password=generate_random_string(10),
        )
        response = self.stub.CreateUser(request)
        logger.info(f"CreateUser: {response}")
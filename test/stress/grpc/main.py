import logging
import time
import string
import grpc
import random

from typing import Any, Callable

import grpc.experimental.gevent as grpc_gevent
from grpc_interceptor import ClientInterceptor
from locust import User, task

from user import user_pb2, user_pb2_grpc


logger = logging.getLogger(__name__)

# patch grpc so that it uses gevent instead of asyncio
grpc_gevent.init_gevent()


class LocustInterceptor(ClientInterceptor):
    def __init__(self, environment, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.env = environment

    def intercept(
        self,
        method: Callable,
        request_or_iterator: Any,
        call_details: grpc.ClientCallDetails,
    ):
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


class GRPCUserService(User):
    host = "localhost:50051"
    stub_class = user_pb2_grpc.UserServiceStub

    def on_start(self):
        super().on_start()
        self._channel = grpc.insecure_channel(self.host)
        interceptor = LocustInterceptor(environment=self.environment)
        self._channel = grpc.intercept_channel(self._channel, interceptor)
        self.stub = user_pb2_grpc.UserServiceStub(self._channel)

    def on_stop(self):
        self._channel.close()
        super().on_stop()

    def _random_string(self, length):
        return "".join(random.choices(string.ascii_letters, k=length))

    @task
    def create_user(self):
        rd_str = self._random_string(10)
        user_input = user_pb2.CreateUserRequest(
            name=rd_str,
            email=rd_str + "@gmail.com",
            password=rd_str,
        )
        self.stub.CreateUser(user_input)

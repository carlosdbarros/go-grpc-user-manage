# test/locust/locust_grpc.py
import logging
import time
import grpc
import grpc._common
from locust import FastHttpUser, User, task, events, between

from user import user_pb2_grpc, user_pb2

from utils import generate_user, generate_user_request, LocustClientInterceptor

import grpc.experimental.gevent as grpc_gevent

# patch grpc so that it uses gevent instead of asyncio
grpc_gevent.init_gevent()
logger = logging.getLogger(__name__)


class GrpcUser(User):
    abstract = True
    min_wait = 1000
    max_wait = 10000

    host = None
    stub_class = None

    def __init__(self, environment, *args, **kwargs):
        super().__init__(environment, *args, **kwargs)
        for attr in ["host", "stub_class"]:
            if not getattr(self, attr):
                raise AttributeError(
                    f"{attr} is required for {self.__class__.__name__}"
                )

        self._channel = grpc.insecure_channel(self.host.lstrip("http://"))
        grpc.channel_ready_future(self._channel).result(timeout=15)
        self.stub = self.stub_class(self._channel)
        interceptor = LocustClientInterceptor(environment=environment)
        self._channel_intercept = grpc.intercept_channel(self._channel, interceptor)
        self.stub_intercept = self.stub_class(self._channel_intercept)

    def on_stop(self):
        self._channel_intercept.close()
        self._channel.close()
        super().on_stop()


def generate_user_iterator(size=1):
    for _ in range(size):
        yield generate_user_request(generate_user())


BATCH_SIZE = 10

class GrpcCreateUser(GrpcUser):
    host = "http://localhost:50051"
    stub_class = user_pb2_grpc.UserServiceStub

    @task
    def create_user_unary(self):
        self.stub_intercept.CreateUser(generate_user_request(generate_user()))

    @task
    def create_user_unary_without_interceptor(self):
        fn_name = "create_user_unary_without_interceptor"
        response = None
        exception = None
        response_size = 0
        start_time = time.perf_counter()
        try:
            response = self.stub.CreateUser(generate_user_request(generate_user()))
            response_size = response.ByteSize()
        except grpc.RpcError as e:
            exception = e
        self.environment.events.request.fire(
            request_type="gRPC",
            name=fn_name,
            response_time=(time.perf_counter() - start_time) * 1000,
            response=response,
            response_length=response_size,
            exception=exception,
        )

    @task
    def create_user_stream_stream(self):
        self.stub_intercept.CreateUserStream(generate_user_iterator())

    @task
    def create_user_unary_batch(self):
        for _ in range(BATCH_SIZE):
            self.stub_intercept.CreateUser(generate_user_request(generate_user()))

    @task
    def create_user_stream_stream_batch(self):
        for _ in self.stub_intercept.CreateUserStream(generate_user_iterator(BATCH_SIZE)):
            pass

    @task
    def create_user_stream_stream_batch_without_interceptor(self):
        exception = None
        response_size = 0
        start_time = time.perf_counter()
        for response in self.stub.CreateUserStream(generate_user_iterator(BATCH_SIZE)):
            logger.info(response)
            self.environment.events.request.fire(
                request_type="gRPC",
                name="create_user_stream_stream_batch",
                response_time=(time.perf_counter() - start_time) * 1000,
                response=response,
                response_length=response_size,
                exception=exception,
            )
            start_time = time.perf_counter()

class HttpCreateUser(FastHttpUser):
    min_wait = 1000
    max_wait = 10000
    host = "http://localhost:8080"

    @task
    def create_user(self):
        self.client.post("/users", json=generate_user())


    @task
    def create_user_batch(self):
        for _ in range(BATCH_SIZE):
            self.client.post("/users", json=generate_user())

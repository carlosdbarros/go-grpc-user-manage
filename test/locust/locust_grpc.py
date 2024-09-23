# test/locust/locust_grpc.py
import logging
import time
import grpc
import grpc._common
from locust import User, task, events, tag

from user import user_pb2_grpc

from utils import generate_user, generate_user_request, LocustInterceptor

import grpc.experimental.gevent as grpc_gevent

# patch grpc so that it uses gevent instead of asyncio
grpc_gevent.init_gevent()
logger = logging.getLogger(__name__)


class GrpcUser(User):
    abstract = True
    min_wait = 1000
    max_wait = 10000
    # wait_time = between(0.1, 0.5)

    host = None
    stub_class = None

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        for attr in ["host", "stub_class"]:
            if not getattr(self, attr):
                raise AttributeError(f"{attr} is required for {self.__class__.__name__}")

    def on_start(self):
        super().on_start()
        logger.info("GrpcUser => on_start")
        self._channel = grpc.insecure_channel(self.host.lstrip("http://"))
        grpc.channel_ready_future(self._channel).result(timeout=15)
        self._channel = grpc.intercept_channel(self._channel, LocustInterceptor(self.environment))
        self.stub = self.stub_class(self._channel)

    def on_stop(self):
        logger.info("GrpcUser => on_stop")
        self._channel.close()
        super().on_stop()


def generate_user_iterator(size=1):
    for _ in range(size):
        yield generate_user_request(generate_user())


class GrpcCreateUser(GrpcUser):
    host = "http://localhost:50051"
    stub_class = user_pb2_grpc.UserServiceStub

    BATCH_SIZE = 10
    WAIT_TIME = 0.01

    # @task
    def create_user_unary(self):
        function_name = self.create_user_unary.__name__
        input_data = generate_user_request(generate_user())
        start_time = time.perf_counter()
        response_byte_size = 0
        try:
            response_message = self.stub.CreateUser(input_data)
            response_byte_size = response_message.ByteSize()
        except Exception as e:
            events.request.fire(
                request_type="gRPC",
                name=function_name,
                response_time=(time.perf_counter() - start_time) * 1000,
                exception=e,
                response_length=response_byte_size,
            )
        else:
            events.request.fire(
                request_type="grpc",
                name=function_name,
                response_time=(time.perf_counter() - start_time) * 1000,
                response_length=response_byte_size,
            )

    # @task
    def create_user_unary_batch(self):
        function_name = self.create_user_unary_batch.__name__
        start_time = time.perf_counter()
        response_byte_size = 0
        for _ in range(self.BATCH_SIZE):
            try:
                input_data = generate_user_request(generate_user())
                response_message = self.stub.CreateUser(input_data)
                response_byte_size = response_message.ByteSize()
            except Exception as e:
                events.request.fire(
                    request_type="gRPC",
                    name=function_name,
                    response_time=(time.perf_counter() - start_time) * 1000,
                    exception=e,
                    response_length=response_byte_size,
                )
            else:
                events.request.fire(
                    request_type="gRPC",
                    name=function_name,
                    response_time=(time.perf_counter() - start_time) * 1000,
                    response_length=response_byte_size,
                )

    @task
    def create_user_stream_stream(self):
        input_data = generate_user_iterator()
        for res in self.stub.CreateUserStream(input_data):
            logger.info(f"Response: {res}")
        # function_name = self.create_user_stream_stream.__name__
        # full_start_time = time.perf_counter()
        # try:
        #     start_time = time.perf_counter()
        #     input_data = generate_user_iterator()
        #     for _ in self.stub.CreateUserStream(input_data):
        #         logger.info(f"GrpcCreateUser => create_user_stream_stream => Response: {response}")
        #         events.request.fire(
        #             request_type="gRPC",
        #             name=function_name,
        #             response_time=(time.perf_counter() - start_time) * 1000,
        #             response_length=response.ByteSize(),
        #             response=response,
        #         )
        #         start_time = time.perf_counter()
        #         time.sleep(self.WAIT_TIME)
        # except Exception as e:
        #     events.request.fire(
        #         request_type="gRPC",
        #         name=function_name,
        #         response_time=(time.perf_counter() - full_start_time) * 1000,
        #         response_length=0,
        #         exception=e,
        #     )
        #     pass

    # @task
    def create_user_stream_stream_batch(self):
        function_name = self.create_user_stream_stream_batch.__name__
        full_start_time = time.perf_counter()
        try:
            start_time = time.perf_counter()
            input_data = generate_user_iterator(self.BATCH_SIZE)
            for response in self.stub.CreateUserStream(input_data):
                events.request.fire(
                    request_type="gRPC",
                    name=function_name,
                    response_time=(time.perf_counter() - start_time) * 1000,
                    response_length=response.ByteSize(),
                    response=response,
                )
                start_time = time.perf_counter()
                time.sleep(self.WAIT_TIME)
        except Exception as e:
            logger.error(e)
            events.request.fire(
                request_type="gRPC",
                name=function_name,
                response_time=(time.perf_counter() - full_start_time) * 1000,
                response_length=0,
                response=None,
                exception=e,
            )

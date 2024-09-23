import logging
import time
import grpc
import faker
import random
import typing as t

from locust.env import Environment

from user import user_pb2

faker = faker.Faker()
logger = logging.getLogger(__name__)


users = [
    {
        "name": faker.name(),
        "email": faker.email(),
        "password": faker.password(),
    }
    for _ in range(10)
]


def generate_user():
    return {
        "name": "Test User",
        "email": "t@t.com",
        "password": "123"
    }


def generate_user_request(user):
    return user_pb2.CreateUserRequest(
        name=user["name"],
        email=user["email"],
        password=user["password"],
    )



class LocustClientInterceptor(
    grpc.UnaryUnaryClientInterceptor, grpc.StreamStreamClientInterceptor
):
    def __init__(self, environment: Environment, *args, **kwargs):
        self.env = environment

    def intercept_call(self, continuation, client_call_details, request_or_iterator):
        future = None
        response = None
        exception = None
        response_size = 0
        start_time = time.perf_counter()
        try:
            future = continuation(client_call_details, request_or_iterator)
        except grpc.RpcError as e:
            exception = e
        else:
            unary_request_types = (user_pb2.CreateUserRequest,)
            if isinstance(request_or_iterator, unary_request_types):
                response = future.result()
                response_size = response.ByteSize()

        self.env.events.request.fire(
            request_type="gRPC",
            name=client_call_details.method,
            response_time=(time.perf_counter() - start_time) * 1000,
            response=response,
            response_length=response_size,
            exception=exception,
        )
        return future

    def intercept_unary_unary(self, continuation, client_call_details, request):
        return self.intercept_call(continuation, client_call_details, request)

    def intercept_stream_stream(self, continuation, client_call_details, request_iterator):
        return self.intercept_call(continuation, client_call_details, request_iterator)

import logging
import time
import grpc
import faker
import random

from user import user_pb2

from typing import Any, Callable, Iterator

import grpc_interceptor
from grpc_interceptor.client import ClientInterceptorReturnType


faker = faker.Faker()
logger = logging.getLogger(__name__)


users = [
    {
        "name": faker.name(),
        "email": faker.email(),
        "password": faker.password(),
    }
    for _ in range(30)
]


def generate_user():
    return random.choice(users)


def generate_user_request(user):
    return user_pb2.CreateUserRequest(
        name=user["name"],
        email=user["email"],
        password=user["password"],
    )


class LocustInterceptor(grpc_interceptor.ClientInterceptor):
    def __init__(self, environment, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.env = environment

    def intercept(
        self,
        method: Callable[..., Any],
        request_or_iterator: Any,
        call_details: grpc.ClientCallDetails,
    ) -> ClientInterceptorReturnType:
        response = None
        response_length = 0
        start_time = time.perf_counter()
        try:
            response = method(request_or_iterator, call_details)
            response_length = response.result().ByteSize()
        except grpc.RpcError as e:
            logger.error(f"LocustInterceptor => Error: {e}")
            self.env.events.request.fire(
                request_type="gRPC",
                name=call_details.method,
                response_time=(time.perf_counter() - start_time) * 1000,
                exception=e,
                response_length=response_length,
            )

        self.env.events.request.fire(
            request_type="gRPC",
            name=call_details.method,
            response_time=(time.perf_counter() - start_time) * 1000,
            response_length=response_length,
            response=response,
            context=None,
        )
        return response

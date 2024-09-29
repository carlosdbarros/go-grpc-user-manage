import time
import string
import logging
import faker
import random
import typing as t
from gevent import monkey
from functools import wraps

import grpc
import grpc.experimental.gevent as grpc_gevent

from locust import User
from locust.env import Environment
from locust.exception import LocustError

from user import user_pb2

# patch grpc so that it uses gevent instead of asyncio
monkey.patch_all()
grpc_gevent.init_gevent()

faker = faker.Faker()
logger = logging.getLogger(__name__)

STR_SIZE = 300
LIST_SIZE = 30
USERS_SIZE = 1
USE_FAKER = False

STATIC_STR = "".join(random.choices(string.ascii_letters, k=STR_SIZE))
STATIC_EMAIL = f"{STATIC_STR}@teste.com"

users = [
    {
        "name": faker.name() if USE_FAKER else STATIC_STR,
        "email": faker.email() if USE_FAKER else STATIC_EMAIL,
        "password": STATIC_STR,
    }
    for _ in range(USERS_SIZE)
]

users_addresses = [
    {
        "name": faker.name() if USE_FAKER else STATIC_STR,
        "emails": [
            faker.email() if USE_FAKER else STATIC_EMAIL for _ in range(LIST_SIZE)
        ],
        "phones": [
            faker.phone_number() if USE_FAKER else STATIC_STR for _ in range(LIST_SIZE)
        ],
        "addresses": [
            {
                "street": faker.street_name() if USE_FAKER else STATIC_STR,
                "number": faker.building_number() if USE_FAKER else STATIC_STR,
                "complement": faker.secondary_address() if USE_FAKER else STATIC_STR,
                "city": faker.city() if USE_FAKER else STATIC_STR,
                "state": faker.state() if USE_FAKER else STATIC_STR,
                "country": faker.country() if USE_FAKER else STATIC_STR,
                "zipCode": faker.zipcode() if USE_FAKER else STATIC_STR,
            }
            for _ in range(LIST_SIZE)
        ],
    }
    for _ in range(USERS_SIZE)
]


def generate_user_address():
    return users_addresses[0]
    # return random.choice(users_addresses)


def generate_user_address_request(user_address=None):
    user_address = user_address or generate_user_address()
    return user_pb2.CreateUserAddressRequest(**user_address)


def generate_user():
    return users[0]
    # return random.choice(users)


def generate_user_request(user=None):
    user = user or generate_user()
    return user_pb2.CreateUserRequest(**user)


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
        total_time = 0
        try:
            future = continuation(client_call_details, request_or_iterator)
            total_time = (time.perf_counter() - start_time) * 1000
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
            response_time=total_time,
            response=response,
            response_length=response_size,
            exception=exception,
        )
        return future

    def intercept_unary_unary(self, continuation, client_call_details, request):
        return self.intercept_call(continuation, client_call_details, request)

    def intercept_stream_stream(
            self, continuation, client_call_details, request_iterator
    ):
        return self.intercept_call(continuation, client_call_details, request_iterator)


class LocustHttpUser(User):
    abstract = True
    host = "http://localhost:8080"

    def on_stop(self, force=False):
        super().stop(force=True)


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
            except Exception as e:
                logger.error(f"Error in {fn.__name__}: {e}")

            response_time = (time.perf_counter() - start_perf_counter) * 1000
            response_size = response.ByteSize() if response else 0
            name = rpc_name.split("/")[-1] if rpc_name else fn.__name__
            if not response:
                logger.info(f"Response is None in {fn.__name__}")

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


class GrpcUser(User):
    abstract = True
    host = "localhost:50051"
    channel: grpc.Channel = None
    compression = None
    stub_class: t.Any = None
    stub: t.Any = None

    def __init__(self, environment):
        super().__init__(environment)
        self.channel = self.__class__.make_channel(self.host or self.environment.host)
        self.stub = self.__class__.make_stub(self.channel)
        self.channel_closed = False

    def on_stop(self, force=False):
        self.channel_closed = True
        time.sleep(.5)
        self.channel.close()
        super().stop(force=True)

    @classmethod
    def make_channel(cls, host: str) -> grpc.Channel:
        if "https" in host:
            raise LocustError("There is no implementation for secure gRPC yet")
        target = host.lstrip("http://")
        channel = grpc.insecure_channel(target, compression=cls.compression)
        grpc.channel_ready_future(channel).result(timeout=10)
        return channel

    @classmethod
    def make_stub(cls, channel: grpc.Channel):
        return cls.stub_class(channel)

import time
import string
import logging
import grpc
import faker
import random
import typing as t

from locust.env import Environment

from user import user_pb2

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

usersAddresses = [
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
    return random.choice(usersAddresses)
    # return {
    #     "name": STATIC_STR,
    #     "emails": [ f"{STATIC_STR}@teste.com" for _ in range(LIST_SIZE)],
    #     "phones": [ STATIC_STR for _ in range(LIST_SIZE)],
    #     "addresses": [
    #         {
    #             "street": STATIC_STR,
    #             "number": STATIC_STR,
    #             "complement": STATIC_STR,
    #             "city": STATIC_STR,
    #             "state": STATIC_STR,
    #             "country": STATIC_STR,
    #             "zipCode": STATIC_STR,
    #         }
    #         for _ in range(LIST_SIZE)
    #     ],
    # }


def generate_user_address_request(user_address):
    return user_pb2.CreateUserAddressRequest(**user_address)


def generate_user():
    return random.choice(users)
    # return {
    #     "name": STATIC_STR,
    #     "email": f"{STATIC_STR}@teste.com",
    #     "password": "123456",
    # }


def generate_user_request(user):
    return user_pb2.CreateUserRequest(**user)


def generate_user_iterator(size=1):
    for _ in range(size):
        yield generate_user_request(generate_user())


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

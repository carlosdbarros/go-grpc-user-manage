import queue
import grpc.experimental.gevent as grpc_gevent

from gevent import monkey

from locust import FastHttpUser, HttpUser, task

from utils import (
    GrpcUser,
    LocustHttpUser,
    make_user,
    make_user_request,
    make_user_address,
    make_user_address_request,
    grpc_stopwatch,
)


# patch grpc so that it uses gevent instead of asyncio
monkey.patch_all()
grpc_gevent.init_gevent()


class GRPCLoadTest(GrpcUser):

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUser")
    def create_user_unary(self):
        return self.stub.CreateUser(make_user_request())

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserAddress")
    def create_user_address_unary(self):
        return self.stub.CreateUserAddress(make_user_address_request())


class GRPCStreamLoadTest(GrpcUser):
    user_queue = queue.Queue()
    user_address_queue = queue.Queue()
    sender_user = None
    sender_user_address = None

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.sender_user = self.stub.CreateUserBidirectional(
            iter(self.user_queue.get, None),
        )
        self.sender_user_address = self.stub.CreateUserAddressBidirectional(
            iter(self.user_address_queue.get, None),
        )

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserBidirectional")
    def create_user_bidirectional_stream(self):
        self.user_queue.put(make_user_request())
        return next(self.sender_user)

    @task
    @grpc_stopwatch("grpc.user.UserService/CreateUserAddressBidirectional")
    def create_user_address_bidirectional_stream(self):
        self.user_address_queue.put(make_user_address_request())
        return next(self.sender_user_address)


class FastHttpLoadTest(LocustHttpUser, FastHttpUser):

    @task
    def create_user(self):
        self.client.post("/users", json=make_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=make_user_address())


class HttpLoadTest(LocustHttpUser, HttpUser):

    @task
    def create_user(self):
        self.client.post("/users", json=make_user())

    @task
    def create_user_address(self):
        self.client.post("/users-address", json=make_user_address())

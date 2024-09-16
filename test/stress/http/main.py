import logging
import string
import random

from locust import HttpUser, task


logger = logging.getLogger(__name__)


class HTTPUSerService(HttpUser):
    host = "http://localhost:8080"

    def on_start(self):
        self.client.headers["Content-Type"] = "application/json"
        self.client.headers["Accept"] = "application/json"

    def on_stop(self):
        self.client.close()

    def _random_string(self, length):
        return "".join(random.choices(string.ascii_letters, k=length))

    @task
    def create_user(self):
        rd_str = self._random_string(10)
        user_input = {
            "name": rd_str,
            "email": f"{rd_str}@{rd_str}.com",
            "password": rd_str,
        }
        self.client.post(
            "/users",
            json=user_input,
        )

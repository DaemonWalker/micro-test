from locust import HttpUser, task

class MtApp(HttpUser):
    @task
    def hello_world(self):
        self.client.get("/")
        self.client.post("/buy", json={    "user": "1@qq.com",    "itemId": 1,    "count": 1})
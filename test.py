import requests
import json

for i in range(1):
    data = json.dumps({"id": i})
    with requests.post("http://localhost:8888/test/form", data=data) as r:
        print(r.status_code)
        print(r.text)
        pass
import json
import requests

print("POST NEW CUSTOMER")
req = requests.post("http://localhost:8888/management/customer", params={"apikey":"7q1qcqmsxnvw"})
print(req.json())


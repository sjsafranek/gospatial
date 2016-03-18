import json
import requests

print()
print("POST NEW CUSTOMER")
req = requests.post("http://localhost:8888/management/customer", params={"apikey":"7q1qcqmsxnvw"})
res = json.loads(req.json())
apikey = res['apikey']


req = requests.post("http://localhost:8888/api/v1/layer", params={"apikey": apikey})
res = json.loads(req.json())
ds = res["datasource"]
print(res)

req = requests.get("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": apikey})
res = req.json()
print(res)


print()
print("POST FEATURE")
payload = {
	"geometry": {
		"type": "Point",
		"coordinates": [10,-10]
	},
	"properties": {
		"name": "test point 1"
	}
}
req = requests.post("http://localhost:8888/api/v1/layer/" + ds + "/feature", params={"apikey": apikey}, data=json.dumps(payload))
print(req.json())


print()
print("GET FEATURE")
req = requests.get("http://localhost:8888/api/v1/layer/" + ds + "/feature/0", params={"apikey": apikey})
print(req.json())


print()
print("DELETE FEATURE")
req = requests.delete("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": apikey})
print(req.json())

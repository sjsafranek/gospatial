import requests

req = requests.put(
	# "http://52.18.218.56:8888/api/v1/layer/73468ae5bc564249a782475d01a0b0ac",
	"https://gis.internalpositioning.com/api/v1/layer/73468ae5bc564249a782475d01a0b0ac",
	params = {
		"apikey": "PuQeSdi8oSCb",
		"authkey": "7q1qcqmsxnvw"
	}
)

https://gis.internalpositioning.com/api/v1/layers?apikey=PuQeSdi8oSCb

print(req.text)

# http://52.18.218.56:8888/api/v1/layer/bd90880b76554d9da9fd31fab72c30d9?apikey=5wbhavVYjun9



req = requests.put( "https://gis.internalpositioning.com/api/v1/layer/73468ae5bc564249a782475d01a0b0ac",
	params = { "apikey": "PuQeSdi8oSCb","authkey": "7q1qcqmsxnvw"}
)

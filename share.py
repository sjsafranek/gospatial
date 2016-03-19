import requests

req = requests.put(
	"http://52.18.218.56:8888/api/v1/layer/bd90880b76554d9da9fd31fab72c30d9",
	params = {
		#"apikey": "I2lus9Vv4rSj",
		"apikey": "5wbhavVYjun9",
		"authkey": "7q1qcqmsxnvw"
	}
)

print(req.text)

# http://52.18.218.56:8888/api/v1/layer/bd90880b76554d9da9fd31fab72c30d9?apikey=5wbhavVYjun9
#!/bin/python
# -*- coding: utf-8 -*-
import json
import argparse
import requests
import unittest

class GoSpatialTester(unittest.TestCase):

	# @classmethod
	# def setUpClass(cls):
	# 	cls.myCustomer = myCustomer

	# @classmethod
	# def tearDownClass(cls):
	# 	print("delete")

	def test_api(self):

		print("[POST] NEW CUSTOMER")
		req = requests.post("http://localhost:8888/management/customer", params={"auth":"7q1qcqmsxnvw"})
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		apikey = res['apikey']
		print("[ok]", apikey)

		print("[POST] NEW DATASOURCE")
		req = requests.post("http://localhost:8888/api/v1/layer", params={"apikey": apikey})
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		ds = res["datasource"]
		print("[ok]", ds)

		print("[GET] DATASOURCE LIST")
		req = requests.get("http://localhost:8888/api/v1/layers", params={"apikey": apikey})
		self.assertEqual(200, req.status_code)
		res = req.json()
		self.assertEqual(apikey, res['apikey'])
		self.assertEqual(ds, res['datasources'][0])
		print("[ok]", ds)

		print("[GET] DATASOURCE")
		req = requests.get("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": apikey})
		self.assertEqual(200, req.status_code)

		print("[POST] NEW FEATURE")
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
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		self.assertEqual(ds, res["datasource"])
		print("[ok]", res["datasource"])

		print("[GET] FEATURE")
		req = requests.get("http://localhost:8888/api/v1/layer/" + ds + "/feature/0", params={"apikey": apikey})
		self.assertEqual(200, req.status_code)
		print("[ok]", ds, 0)

		print("[DELETE] Layer")
		req = requests.delete("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": apikey})
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		self.assertEqual(ds, res["datasource"])
		print("[ok]", res["datasource"])


if __name__ == "__main__":
	parser = argparse.ArgumentParser(description='GoSpatial Unittester')
	parser.add_argument('-c', type=str, help='customer apikey')
	args = parser.parse_args()
	unittest.main(warnings='ignore')


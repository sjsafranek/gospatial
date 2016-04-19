#!/bin/python
# -*- coding: utf-8 -*-
import sys
import json
import argparse
import requests
import unittest
import websocket
import logging


logger = logging.getLogger()
logger.level = logging.DEBUG
stream_handler = logging.StreamHandler(sys.stdout)
stream_handler.setFormatter(logging.Formatter("%(asctime)s [%(threadName)-12.12s] [%(levelname)-5.5s]  %(message)s"))
logger.addHandler(stream_handler)


class GoSpatialTester(unittest.TestCase):

	@classmethod
	def setUpClass(cls):
		logging.info("[POST] NEW CUSTOMER")
		req = requests.post("http://localhost:8888/management/customer", params={"authkey":"7q1qcqmsxnvw"})
		logger.info(req.text)
		if req.status_code is not 200:
			logging.error("Could not create apikey")
			raise ValueError("Could not create apikey")
		res = json.loads(req.json())
		apikey = res['apikey']
		logging.info("[ok] " + apikey)
		cls.apikey = apikey

	# @classmethod
	# def tearDownClass(cls):
	# 	print("delete")

	def test_api(self):

		logging.info("[POST] NEW DATASOURCE")
		req = requests.post("http://localhost:8888/api/v1/layer", params={"apikey": self.apikey})
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		ds = res["datasource"]
		logging.info("[ok] " + ds)

		logging.info("[GET] DATASOURCE LIST")
		req = requests.get("http://localhost:8888/api/v1/layers", params={"apikey": self.apikey})
		self.assertEqual(200, req.status_code)
		res = req.json()
		self.assertEqual(self.apikey, res['apikey'])
		self.assertEqual(ds, res['datasources'][0])
		logging.info("[ok] " + ds)

		logging.info("[GET] DATASOURCE")
		req = requests.get("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": self.apikey})
		self.assertEqual(200, req.status_code)

		logging.info("[POST] NEW FEATURE")
		payload = {
			"geometry": {
				"type": "Point",
				"coordinates": [10,-10]
			},
			"properties": {
				"name": "test point 1"
			}
		}
		req = requests.post("http://localhost:8888/api/v1/layer/" + ds + "/feature", params={"apikey": self.apikey}, data=json.dumps(payload))
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		self.assertEqual(ds, res["datasource"])
		logging.info("[ok] " + res["datasource"])

		logging.info("[GET] FEATURE")
		req = requests.get("http://localhost:8888/api/v1/layer/" + ds + "/feature/0", params={"apikey": self.apikey})
		self.assertEqual(200, req.status_code)
		logging.info("[ok] " + ds + " " + str(0))

		# print("[DELETE] Layer")
		# req = requests.delete("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": self.apikey})
		# self.assertEqual(200, req.status_code)
		# res = json.loads(req.json())
		# self.assertEqual(ds, res["datasource"])
		# print("[ok]", res["datasource"])


	def test_websockets(self):
		logging.info("[POST] NEW DATASOURCE")
		req = requests.post("http://localhost:8888/api/v1/layer", params={"apikey": self.apikey})
		self.assertEqual(200, req.status_code)
		res = json.loads(req.json())
		ds = res["datasource"]
		logging.info("[ok] " + ds)

		ws = websocket.create_connection("ws://localhost:8888/ws/" + ds)
		ws.send("sup!")
		result =  ws.recv()
		logging.info(result)


if __name__ == "__main__":
	parser = argparse.ArgumentParser(description='GoSpatial Unittester')
	parser.add_argument('-c', type=str, help='customer apikey')
	parser.add_argument('-p', type=str, help='port')
	args = parser.parse_args()
	unittest.main(warnings='ignore')


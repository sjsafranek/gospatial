#!/bin/python3
# -*- coding: utf-8 -*-

import json
import argparse
import requests
import logging

from Config import Config


class ApiClient(Config):

	def __init__(self, config_file='config.ini'):
		super().__init__(config_file)

	def createDatasource(self):
		logging.info("[POST] NEW DATASOURCE")
		url = "{}/{}".format(self.baseUrl(), 'api/v1/layer')
		resp = requests.post(url, params={"apikey": self.apikey()})
		if 200 == resp.status_code:
			return resp.json()['data']
		else:
			ValueError("Api error: " + resp.text)

	def getDatasources(self):
		logging.info("[GET] DATASOURCE LIST")
		url = "{}/{}".format(self.baseUrl(), 'api/v1/layer')
		resp = requests.get(url, params={"apikey": self.apikey()})
		if 200 == resp.status_code:
			return resp.json()['data']
		else:
			ValueError("Api error: " + resp.text)

	def getDatasource(self, datasource_id):
		logging.info("[GET] DATASOURCE")
		url = "{}/{}/{}".format(self.baseUrl(), 'api/v1/layer', datasource_id)
		resp = requests.get(url, params={"apikey": self.apikey()})
		if 200 == resp.status_code:
			return
		else:
			ValueError("Api error: " + resp.text)

	def createFeature(self, datasource_id):
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
		url = "{}/{}/{}/feature".format(self.baseUrl(), 'api/v1/layer', datasource_id)
		resp = requests.post(url, params={"apikey": self.apikey()}, data=json.dumps(payload))
		if 200 == resp.status_code:
			return
		else:
			ValueError("Api error: " + resp.text)

	def getFeature(self, datasource_id, feature_id):
		logging.info("[GET] FEATURE")
		url = "{}/{}/{}/feature/{}".format(self.baseUrl(), 'api/v1/layer', datasource_id, feature_id)
		resp = requests.get(url, params={"apikey": self.apikey()})
		if 200 == resp.status_code:
			return resp.json()
		else:
			ValueError("Api error: " + resp.text)

		# print("[DELETE] Layer")
		# req = requests.delete("http://localhost:8888/api/v1/layer/" + ds, params={"apikey": self.apikey})
		# self.assertEqual(200, req.status_code)
		# res = json.loads(req.json())
		# self.assertEqual(ds, res["datasource"])
		# print("[ok]", res["datasource"])


'''
if __name__ == "__main__":
	parser = argparse.ArgumentParser(description='GoSpatial Unittester')
	parser.add_argument('-c', type=str, help='customer apikey')
	parser.add_argument('-p', type=str, help='port')
	args = parser.parse_args()
	unittest.main(warnings='ignore')

'''
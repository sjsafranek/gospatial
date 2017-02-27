#!/usr/bin/env python3
import os
import configparser

class Config(object):

	def __init__(self, config_file='config.ini'):
		self._configFile = config_file
		self.config = configparser.ConfigParser()
		if not os.path.exists(self._configFile):
			self._createConfig()
		else:
			self.config.read(self._configFile)

	def _createConfig(self):
		self.config['credentials'] = {}
		self.config['credentials']['url'] = 'http://localhost:8080'
		self.config['credentials']['apikey'] = self._createTestApikey()
		self.config['selenium'] = {}
		self.config['selenium']['driver'] = 'firefox'
		self.config['selenium']['executable'] = 'null'
		with open(self._configFile, 'w') as configfile:
			self.config.write(configfile)

	def _createTestApikey(self):
		import TcpClient
		return TcpClient.newApikey()

	def apikey(self):
		return self.config['credentials']['apikey']

	def baseUrl(self):
		return self.config['credentials']['url']

	def driverType(self):
		return self.config['selenium']['driver']

	def driverExecutable(self):
		return self.config['selenium']['executable']



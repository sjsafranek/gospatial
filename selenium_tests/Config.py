#!/usr/bin/python3
# -*- coding: utf-8 -*-
import os
import json
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
		#self.config['credentials']['url'] = input('testing url: ')
		self.config['credentials']['apikey'] = self._createTestApikey()
		self.config['selenium'] = {}
		self.config['selenium']['driver'] = 'firefox'
		#self.config['selenium']['executable_path'] = 'null'
		self.config['selenium/chrome'] = {}
		self.config['selenium/chrome']['executable_path'] = '/usr/bin/google-chrome'
		self.config['selenium/phantomjs'] = {}
		self.config['selenium/phantomjs']['executable_path'] = 'null'
		self.saveConfig()

	def saveConfig(self):
		with open(self._configFile, 'w') as configfile:
			self.config.write(configfile)

	def _createTestApikey(self):
		import TcpClient
		return TcpClient.newApikey()

	@property
	def apikey(self):
		return self.config['credentials']['apikey']
		# return self.config['credentials'].get('apikey')

	def setApikey(self, apikey, save=False):
		self.config['credentials']['apikey'] = apikey
		if save: 
			self.saveConfig()

	@property
	def baseUrl(self):
		return self.config['credentials']['url']
		# return self.config['credentials'].get('url')

	def setBaseUrl(self, url, save=False):
		self.config['credentials']['url'] = url
		if save: 
			self.saveConfig()

	@property
	def driverType(self):
		return self.config['selenium']['driver']

	@property
	def driverExecutable(self):
		return self.config['selenium/'+self.driverType]['executable_path']

	def setDriverType(self, driver):
		if driver not in ['firefox', 'chrome', 'ie', 'opera', 'phantomjs']:
			ValueError('Unsupported driver: '+driver)
		self.config['selenium']['driver'] = driver

	def setdriverExecutable(self, driver, executable):
		if driver not in ['firefox', 'chrome', 'ie', 'opera', 'phantomjs']:
			ValueError('Unsupported driver: '+driver)
		self.config['selenium/'+driver]['executable'] = executable

	@property
	def driverKwargs(self):
		browsers = [
			'firefox',
			'chrome',
			'ie',
			'opera',
			'phantomjs',
			'remote'
		]

		browser_kwargs = dict((k, {}) for k in browsers)
		for browser in browser_kwargs.keys():
			section = 'selenium/%s' % browser
			if section in self.config.sections():
				browser_kwargs[browser] = dict(self.config[section])
				if "desired_capabilities" in self.config[section]:
					browser_kwargs[browser]['desired_capabilities'] = json.loads(self.config[section]['desired_capabilities'])
		#config_browser = self.config['selenium'].get('driver')
		return browser_kwargs[self.driverType]

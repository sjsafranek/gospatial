#!/usr/bin/python3
# -*- coding: utf-8 -*-
import nose
import unittest
from Db4IoTViewerTester import Db4IoTViewerTester

class MapUnittestBase(unittest.TestCase):
	_multiprocess_can_split_ = True
	
	@classmethod
	def setUpClass(cls):
		cls.client = Db4IoTViewerTester()

	@classmethod
	def tearDownClass(cls):
		cls.client.destroy()

	def setUp(self):
		""" prep for test """
		self.client.login("map")

	def tearDown(self):
		""" clean up test """
		self.client.logout()

if __name__ == '__main__':
	unittest.main(warnings='ignore')

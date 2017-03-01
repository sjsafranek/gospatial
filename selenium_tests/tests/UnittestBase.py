#!/usr/bin/python3
# -*- coding: utf-8 -*-
import nose
import unittest
from WebClient import WebClient

class DashboardUnittestBase(unittest.TestCase):

	_multiprocess_can_split_ = True
	
	@classmethod
	def setUpClass(cls):
		cls.client = WebClient()

	@classmethod
	def tearDownClass(cls):
		cls.client.destroy()

	def setUp(self):
		""" prep for test """
		self.client.dashboardPage()

	#def tearDown(self):
		""" clean up test """
		#self.client.logout()

if __name__ == '__main__':
	unittest.main(warnings='ignore')

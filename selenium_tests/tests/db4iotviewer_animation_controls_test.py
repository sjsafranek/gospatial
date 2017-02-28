#!/usr/bin/python3
# -*- coding: utf-8 -*-

from time import sleep
import unittest
from selenium.webdriver.common.by import By
from Db4IoTViewerTester import Db4IoTViewerTester


class AnimationControlsTest(unittest.TestCase):
	@classmethod
	def setUpClass(cls):
		cls.client = Db4IoTViewerTester()

	@classmethod
	def tearDownClass(cls):
		#cls.client.logout()
		cls.client.destroy()

	def setUp(self):
		""" prep for test """
		self.client.login("map")

	def test_forward_animation(self):
		self.client.setEventTimestampToCenter()
		self.client.play()
		t1 = self.client.getEventTimeStamp()
		sleep(5)
		t2 = self.client.getEventTimeStamp()
		self.assertTrue(t2 > t1)

	def test_backward_animation(self):
		self.client.setEventTimestampToCenter()
		self.client.back()
		t1 = self.client.getEventTimeStamp()
		sleep(5)
		t2 = self.client.getEventTimeStamp()
		self.assertTrue(t2 < t1)
		
	def test_pause_animation(self):
		self.client.setEventTimestampToCenter()
		self.client.play()
		sleep(3)
		self.client.pause()
		t1 = self.client.getEventTimeStamp()
		sleep(4)
		t2 = self.client.getEventTimeStamp()
		self.assertTrue(t2 == t1)
		
	def tearDown(self):
		""" clean up test """
		self.client.logout()


if __name__ == '__main__':
	unittest.main(warnings='ignore')

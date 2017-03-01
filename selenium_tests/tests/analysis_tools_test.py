#!/usr/bin/python3
# -*- coding: utf-8 -*-

import nose
from time import sleep
import unittest
from selenium.webdriver.common.by import By
from Db4IoTViewerTester import Db4IoTViewerTester
from MapUnittest import MapUnittestBase

class AnalsysToolsTest(MapUnittestBase):

	def test_open_analysis_tools(self):
		self.client.setEventTimestampToCenter()
		self.client.openSidepanel()
		# check height
		height_begin = self.client.getElem('#miscFeatureMethods').size['height']
		# open all tools
		self.client.openAnalysisTool('piechart')
		self.client.openAnalysisTool('linechart')
		self.client.openAnalysisTool('barchart')
		self.client.openAnalysisTool('export_csv')
		# check height
		height_open = self.client.getElem('#miscFeatureMethods').size['height']
		# close all tools
		self.client.closeAnalysisTool('piechart')
		self.client.closeAnalysisTool('linechart')
		self.client.closeAnalysisTool('barchart')
		self.client.closeAnalysisTool('export_csv')
		# check height
		height_end = self.client.getElem('#miscFeatureMethods').size['height']
		# compare heights
		self.assertTrue(height_open > height_begin)
		self.assertTrue(height_open > height_end)


if __name__ == '__main__':
	unittest.main(warnings='ignore')

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
		cls.client.destroy()

	def setUp(self):
		""" prep for test """
		self.client.login("map")

	def test_open_close_sidepanel(self):
		self.client.openSidepanel()
		self.assertTrue( self.client.getElem('div#sidepanel').is_displayed() )
		self.client.closeSidepanel()
		self.assertFalse( self.client.getElem('div#sidepanel').is_displayed() )

	def test_sidepanel_sections(self):
		# widgets section
		elems = self.client.getElems('#miscFeatureMethods div.panel')
		self.assertTrue( len(elems) > 0 )
		# maplayers section
		elems = self.client.getElems('#animationMethods div.panel')
		self.assertTrue( len(elems) > 0 )
		# filter section
		elems = self.client.getElems('#selectors div.field')
		self.assertTrue( len(elems) > 0 )

	def test_change_maplayers(self):
		self.client.setEventTimestampToCenter()
		self.client.openSidepanel()
		for elem in self.client.getElems(".mapLayer"):
			if elem.is_displayed():
				maplayer = elem.get_attribute('layer_name')
				self.client.setMapLayer(maplayer)
				sleep(1)
				if 'grid' == maplayer:
					sleep(2.5)
					self.assertTrue(len(self.client.getElems('g path.leaflet-clickable')) > 450)
				if 'marker' == maplayer:
					sleep(2.5)
					num_icons = len(self.client.getElems('div.svgIcon'))
					num_shadows = len(self.client.getElems('svg.symbol.shadow'))
					num_labels = len(self.client.getElems('span.markerLabel'))
					self.assertTrue(num_icons == num_shadows)
					self.assertTrue(num_icons == num_labels)
					self.assertTrue(num_shadows == num_labels)
				self.assertTrue(maplayer == self.client.getMapLayer())
	
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
		self.assertTrue(height_end > height_begin)
		self.assertTrue(height_open > height_end)

	def tearDown(self):
		""" clean up test """
		self.client.logout()


if __name__ == '__main__':
	unittest.main(warnings='ignore')

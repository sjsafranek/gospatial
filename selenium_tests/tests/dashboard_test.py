#!/usr/bin/python3
# -*- coding: utf-8 -*-

import nose
from time import sleep
import unittest
from selenium.webdriver.common.by import By
from UnittestBase import DashboardUnittestBase

class DashboardTest(DashboardUnittestBase):

	_multiprocess_can_split_ = False

	def test_create_layer(self):
		datasource_ids_begin = self.client.getDatasourceIds()
		self.client.createNewLayer()
		datasource_ids_end = self.client.getDatasourceIds()
		self.assertTrue( len(datasource_ids_end) > len(datasource_ids_begin) )

	def test_delete_layer(self):
		datasource_ids_begin = self.client.getDatasourceIds()
		n = len(datasource_ids_begin)
		self.client.deleteLayer(datasource_ids_begin[n-1])
		datasource_ids_end = self.client.getDatasourceIds()
		self.assertTrue( len(datasource_ids_end) < len(datasource_ids_begin) )


if __name__ == '__main__':
	unittest.main(warnings='ignore')

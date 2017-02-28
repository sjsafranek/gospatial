#!/usr/bin/env python3

import time

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait 
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains

from Config import Config

class WebClient(Config):

	def __init__(self, config_file='config.ini', ):
		super().__init__(config_file)
		self.driver = None
		if 'firefox' == self.driverType():
			self.driver = webdriver.Firefox()
		elif 'chrome' == self.driverType():
			self.driver = webdriver.Chrome(self.driverExecutable())
		else:
			ValueError("Browser driver is unsupported")

	def _getDriver(self):
		browsers = {
			'firefox': webdriver.Firefox,
			'chrome': webdriver.Chrome,
			'ie': webdriver.Ie,
			'opera': webdriver.Opera,
			'phantomjs': webdriver.PhantomJS,
		}

		browser_kwargs = dict((k, {}) for k in browsers.keys())
		for browser in browser_kwargs.keys():
			section = 'selenium/%s' % browser
			if section in self.config.sections():
				browser_kwargs[browser] = dict(self.config[section])

		config_browser = self.config['selenium'].get('driver')
		if config_browser:
			# Fail if set browser invalid
			self._browser = browsers[config_browser]
			self._browser_kwargs = browser_kwargs[config_browser]
		else:
			# Default to using firefox
			self._browser = browsers['firefox']
			self._browser_kwargs = browser_kwargs['firefox']
		
		return self._browser(**self._browser_kwargs)
		#return self.config["selenium"].get('browser')

	def getPage(self, page, refresh=False):
		if "/"+page not in self.driver.current_url or refresh:
			url = "{}/{}?apikey={}".format(self.baseUrl(), page, self.apikey())
			self.driver.get( url )

	def dashboardPage(self):
		self.getPage("dashboard")

	def mapPage(self):
		self.getPage("map")

	def getElem(self, css_selector):
		WebDriverWait(self.driver, 10).until(
			EC.presence_of_element_located( 
				(By.CSS_SELECTOR, css_selector)
			) 
		)
		return self.driver.find_element(By.CSS_SELECTOR, css_selector)

	def getElems(self, css_selector):
		WebDriverWait(self.driver, 10).until( 
			EC.presence_of_element_located( 
				(By.CSS_SELECTOR, css_selector) 
			) 
		)
		return self.driver.find_elements(By.CSS_SELECTOR, css_selector)

	def clickElem(self, css_selector):
		self.getElem(css_selector).click()

	def confirmSwalPopup(self):
		self.clickElem('button.swal2-confirm.swal2-styled')

	def cancleSwalPopup(self):
		self.clickElem('button.swal2-cancel.swal2-styled')

	def checkSwalApiResponse(self):
		WebDriverWait(self.driver, 10).until( EC.visibility_of_element_located( (By.ID, "modalContentId") ) )
		api_response_message = self.driver.find_element(By.ID, "modalContentId").text
		if "success" not in api_response_message:
			ValueError("Api Error: " + api_response_message)

	def createNewLayer(self):
		# navigate to dashboard page
		self.dashboardPage()
		# close popup
		self.clickElem("button.createLayer")
		time.sleep(0.25)
		self.confirmSwalPopup()
		# check api response
		time.sleep(0.25)
		self.checkSwalApiResponse()
		# close popup
		self.confirmSwalPopup()

	def deleteLayer(self, datasource_id):
		# navigate to dashboard page
		self.dashboardPage()
		# get panel
		time.sleep(0.25)
		for button in self.driver.find_elements(By.CLASS_NAME, "toggleVectorLayerOptions"):
			if datasource_id == button.get_attribute("ds_id"):
				button.click()
				time.sleep(0.25)
				for delete_button in self.driver.find_elements(By.CLASS_NAME, "deleteLayer"):
					if datasource_id == delete_button.get_attribute("ds_id"):
						delete_button.click()
						self.confirmSwalPopup()
						time.sleep(0.25)
						self.checkSwalApiResponse()
						self.confirmSwalPopup()
						break
				return
		# datasource_id not found
		ValueError("Datasource not found: " + datasource_id)


 
wc = WebClient()


# browser.find_elements_by_xpath("//*[@type='submit']")

#wc.dashboardPage()
#wc.createNewLayer()
#wc.deleteLayer("18cf1cc5c0db4589833af8b956dcb631")

wc.mapPage()




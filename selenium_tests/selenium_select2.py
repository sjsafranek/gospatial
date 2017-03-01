# original: https://gist.github.com/benselme/5817536
# forked:   https://gist.github.com/rafaelugolini/d2067a8c8c54026ac029
import time
from selenium.webdriver import ActionChains
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

class Select2(object):
    def __init__(self, element):
        self.browser = element.parent
        self.replaced_element = element
        self.element = element

    def click(self, element=None):
        if element is None:
            element = self.element.find_element(By.CSS_SELECTOR, 'span.select2-selection')
        click_element = ActionChains(self.browser)\
            .click_and_hold(element)\
            .release(element)
        click_element.perform()

    def open(self):
        if not self.is_open:
            self.click()
            WebDriverWait(self.browser, 5).until(
               EC.visibility_of_element_located((By.CSS_SELECTOR, '.select2-dropdown')))

    def close(self):
        if self.is_open:
            self.click()

    def select(self, name):
        self.open()
        item_divs = self.dropdown.find_elements(
			By.CSS_SELECTOR,
            'li.select2-results__option span')
        for field in item_divs:
            if field.text == name:
                self.click(field)
                return True
        return False

    @property
    def is_open(self):
        try:
            if self.element.find_element(By.CSS_SELECTOR, 'span.select2-selection').get_attribute('aria-owns'):
                return True
            else:
                return False
        except:
            return False

    @property
    def dropdown(self):
        self.open()
        elem_id = self.element.find_element(By.CSS_SELECTOR, 'span.select2-selection').get_attribute('aria-owns')
        WebDriverWait(self.browser, 5).until(
           EC.visibility_of_element_located((By.ID, elem_id)))
        return self.browser.find_element(By.ID, elem_id)

    @property
    def items(self):
        self.open()
        item_divs = self.dropdown.find_elements_by_css_selector(
            'li.select2-results__option span')
        return [div for div in item_divs]


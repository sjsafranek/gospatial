#!/bin/python

import logging
import logging.handlers

def log(log_name):
	# Set up a specific logger with our desired output level
	logger = logging.getLogger(log_name)
	handler = logging.FileHandler(log_name+'.log')
	logger.setLevel(logging.DEBUG)
	ch = logging.StreamHandler()
	ch.setLevel(logging.DEBUG)
	# create formatter
	formatter = logging.Formatter("%(asctime)s [%(levelname)s] [%(name)s] %(filename)s line:%(lineno)d : %(message)s")
	# add formatter to ch
	handler.setFormatter(formatter)
	ch.setFormatter(formatter)
	# add ch to logger
	logger.addHandler(ch)
	logger.addHandler(handler) 
	return logger

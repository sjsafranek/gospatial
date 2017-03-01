#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys
import nose
import os
import glob
from nose.plugins import Plugin
from nose.plugins.multiprocess import MultiProcess
from nose.plugins.xunit import Xunit
# http://nose.readthedocs.io/en/latest/doc_tests/test_multiprocess/multiprocess.html

#from nose_xunitmp import XunitMP

if __name__ == "__main__":

    args = [
        "-v", 
        "--with-xunit",
        "--xunit-file=results.xml",
        #"--with-xunitmp",
        #"--xunitmp-file=results.xml",
        "--processes=10", 			# parallel tests
        "--process-timeout=120"		# parallel tests
    ]
	
    if sys.argv[1:]:
        args.extend(sys.argv[1:])
    else:
        args.extend(glob.glob(os.path.join("tests", "*_test.py")))

    nose.run(argv=args, plugins=[MultiProcess(), Xunit()])
    #nose.run(argv=args, plugins=[MultiProcess(), XunitMP()])

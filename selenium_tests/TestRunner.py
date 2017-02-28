#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys
import nose
import os
import glob
from nose.plugins import Plugin

class CsvReport(Plugin):
    name = "csv-report"
    results = "Test,Results,Details\n"

    def write_report(self, filename=None):
        if filename is None:
            filename = 'report.csv'

        with open(filename, "w") as report_file:
            report_file.write(self.results)

    def create_zendesk_ticket(self):
        Zendesk.create_ticket(self.results)

    def addSuccess(self, *args, **kwargs):
        test = args[0]
        self.results += "%s,%s,\n" % (test, "Success")

    def addError(self, *args, **kwargs):
        test, error = args[0:2]
        self.results += "%s,%s,%s\n" % (test, "Error", error[1])

    def addFailure(self, *args, **kwargs):
        test, error = args[0:2]
        self.results += "%s,%s,%s\n" % (test, "Failure", error[1])


if __name__ == "__main__":
    csv_report = CsvReport()

    args = ["-v", "--with-csv-report"]
    if sys.argv[1:]:
        args.extend(sys.argv[1:])
    else:
        args.extend(glob.glob(os.path.join("tests", "*.py")))

    nose.run(argv=args, plugins=[csv_report])

    csv_report.write_report()

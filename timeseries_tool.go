package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
)

import (
	"github.com/sjsafranek/DiffDB/diff_db"
	"github.com/sjsafranek/DiffDB/diff_store"
)

const (
	NAME   = "DiffDB Client"
	BINARY = "diff_db_clie"
)

// RuntimeArgs contains all runtime
// arguments available
var RuntimeArgs struct {
	DatabaseLocation string
	Verbose          bool
}

var (
	diffDb diff_db.DiffDb
)

func errorHandler(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func incorrectUsageError() {
	err := fmt.Errorf("Incorrect usage! Nonsensical argument!")
	errorHandler(err)
}

func successHandler(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func usage() {
	fmt.Printf("%s %s\n\n", NAME, "0.0.2")
	fmt.Printf("Usage:\n\t%s [options...] action key [action_args...]\n\n", BINARY)
	fmt.Println(" * action:\tThe action to preform. Supported action(s): GET, SET, DEL")
	fmt.Println(" * action_args:\tVariadic arguments provided to the requested action. Different actions require different arguments")
	fmt.Println("\n")
}

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")

	// handle command line arguements
	flag.Usage = usage
	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.BoolVar(&RuntimeArgs.Verbose, "verbose", false, "verbose")
	flag.Parse()

	// create database object
	diffDb = diff_db.NewDiffDb(RuntimeArgs.DatabaseLocation)

	// get args
	args := flag.Args()

	// If only one arg
	if 1 == len(args) {
		action := args[0]
		if "KEYS" == action {
			keys, err := diffDb.SelectAll()
			if nil != err {
				errorHandler(err)
			}
			result := fmt.Sprintf("%v", keys)
			successHandler(result)
		}
	}

	if 2 > len(args) {
		incorrectUsageError()
	}

	// command line args
	action := args[0]
	key := args[1]

	switch action {
	// get value of key
	case "GET":
		var ddata diff_store.DiffStore

		data, err := diffDb.Load(key)
		if nil != err {
			errorHandler(err)
		}
		ddata.Decode(data)

		if 2 == len(args) {
			enc, _ := ddata.Encode()
			msg := fmt.Sprintf("%s", enc)
			successHandler(msg)
		}

		if "VALUE" == args[2] {
			successHandler(ddata.GetCurrent())
		}

		if "SNAPSHOTS" == args[2] || "TIMESTAMPS" == args[2] {
			msg := fmt.Sprintf("%v", ddata.GetSnapshots())
			successHandler(msg)
		}

		if 4 == len(args) {

			num, err := strconv.ParseInt(args[3], 10, 64)
			if nil != err {
				errorHandler(err)
			}

			if "TIMESTAMP" == args[2] {
				val, err := ddata.GetPreviousByTimestamp(num)
				if nil != err {
					errorHandler(err)
				}
				successHandler(val)
			}

			if "INDEX" == args[2] {
				val, err := ddata.GetPreviousByIndex(int(num))
				if nil != err {
					errorHandler(err)
				}
				successHandler(val)
			}

		}

		if 5 == len(args) {

			num1, err := strconv.ParseInt(args[3], 10, 64)
			if nil != err {
				errorHandler(err)
			}

			num2, err := strconv.ParseInt(args[4], 10, 64)
			if nil != err {
				errorHandler(err)
			}

			if "RANGE" == args[2] {
				vals, err := ddata.GetPreviousWithinTimestampRange(num1, num2)
				if nil != err {
					errorHandler(err)
				}

				msg := "timestamp,value\n"
				for i := range vals {
					msg += fmt.Sprintf("%v,%s\n", i, vals[i])
				}
				successHandler(msg)
			}
		}

		incorrectUsageError()

	// set new value for key
	case "SET":

		// check for data to set as new value
		if 3 > len(args) {
			incorrectUsageError()
		}

		// load key
		var ddata diff_store.DiffStore
		data, err := diffDb.Load(key)
		if nil != err {
			if err.Error() == "Not found" {
				// create new diffstore if key not found in database
				ddata = diff_store.NewDiffStore(key)
			} else {
				errorHandler(err)
			}
		} else {
			ddata.Decode(data)
		}

		// update diffstore
		ddata.Update(args[2])

		// save to database
		enc, err := ddata.Encode()
		if nil != err {
			errorHandler(err)
		}
		diffDb.Save(ddata.Name, enc)

		// print result
		successHandler(ddata.GetCurrent())

	// delete key
	case "DEL":
		err := diffDb.Remove(key)
		if nil != err {
			errorHandler(err)
		}

	default:
		err := fmt.Errorf("Unsupported action %s, cannot process.", args[0])
		errorHandler(err)
	}

}

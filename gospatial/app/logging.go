package app

import (
	"io"
	// "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Info                    *log.Logger
	Debug                   *log.Logger
	Warning                 *log.Logger
	Error                   *log.Logger
	serverLoggerWriter      io.Writer
	networkLoggerWriter     io.Writer
	networkLoggerInfo     *log.Logger
	networkLoggerWarning  *log.Logger
	networkLoggerError    *log.Logger
	networkLoggerInfoIn  *log.Logger
	networkLoggerInfoOut *log.Logger
)

// func networkLoggerInit() {
// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// 	if err != nil {
// 		Error.Fatal(err)
// 	}
// 	log_file := strings.Replace(dir, "bin", "log/network.log", -1)
// 	networkLoggerWriter, err = os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		Error.Fatal("Error opening file: %v", err)
// 	}
// 	networkLoggerInfoIn = log.New(networkLoggerWriter, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerInfoOut = log.New(networkLoggerWriter, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerInfo = log.New(networkLoggerWriter, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerWarning = log.New(networkLoggerWriter, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerError = log.New(networkLoggerWriter, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// }

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}

	// server logging
	serverLogFile := strings.Replace(dir, "bin", "log/server.log", -1)
	serverLoggerWriter, err := os.OpenFile(serverLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer serverLoggerWriter.Close()
	Info = log.New(serverLoggerWriter, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(serverLoggerWriter, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(serverLoggerWriter, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(serverLoggerWriter, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)

	// network logging
	networkLogFile := strings.Replace(dir, "bin", "log/network.log", -1)
	networkLoggerWriter, err = os.OpenFile(networkLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	networkLoggerInfoIn = log.New(networkLoggerWriter, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	networkLoggerInfoOut = log.New(networkLoggerWriter, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	networkLoggerInfo = log.New(networkLoggerWriter, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	networkLoggerWarning = log.New(networkLoggerWriter, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	networkLoggerError = log.New(networkLoggerWriter, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func testLoggerInit() {
	serverLoggerWriter, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer serverLoggerWriter.Close()
	Info = log.New(serverLoggerWriter, "[TESTING] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(serverLoggerWriter, "[TESTING] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(serverLoggerWriter, "[TESTING] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(serverLoggerWriter, "[TESTING] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

// func StdOutMode() {
// 	AppMode = "standard"
// 	Info = log.New(os.Stdout, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Debug = log.New(os.Stdout, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Warning = log.New(os.Stdout, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	Error = log.New(os.Stderr, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerInfo = log.New(os.Stdout, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerWarning = log.New(os.Stdout, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerError = log.New(os.Stderr, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerInfoIn = log.New(os.Stdout, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	networkLoggerInfoOut = log.New(os.Stdout, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// }

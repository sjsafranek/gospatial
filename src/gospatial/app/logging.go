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
	// Trace                  *log.Logger
	Info                    *log.Logger
	Debug                   *log.Logger
	Warning                 *log.Logger
	Error                   *log.Logger
	LogFileHandler          io.Writer
	network_logger_writer   io.Writer
	network_logger_Info     *log.Logger
	network_logger_Warning  *log.Logger
	network_logger_Error    *log.Logger
	network_logger_Info_In  *log.Logger
	network_logger_Info_Out *log.Logger
)

func Network_logger_init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	log_file := strings.Replace(dir, "bin", "log/network.log", -1)
	network_logger_writer, err = os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	network_logger_Info_In = log.New(network_logger_writer, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info_Out = log.New(network_logger_writer, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info = log.New(network_logger_writer, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Warning = log.New(network_logger_writer, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Error = log.New(network_logger_writer, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func init() {
	AppMode = "logging"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	log_file := strings.Replace(dir, "bin", "log/server.log", -1)
	LogFileHandler, err := os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer LogFileHandler.Close()
	// Trace = log.New(LogFileHandler, "TRACE [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(LogFileHandler, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(LogFileHandler, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(LogFileHandler, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(LogFileHandler, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func test_logger_init() {
	AppMode = "testing"
	LogFileHandler, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer LogFileHandler.Close()
	// Trace = log.New(LogFileHandler, "[TESTING] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(LogFileHandler, "[TESTING] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(LogFileHandler, "[TESTING] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(LogFileHandler, "[TESTING] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(LogFileHandler, "[TESTING] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func StdOutMode() {
	AppMode = "standard"
	// Trace = log.New(os.Stdout, "TRACE [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(os.Stdout, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(os.Stdout, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(os.Stderr, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info = log.New(os.Stdout, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Warning = log.New(os.Stdout, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Error = log.New(os.Stderr, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info_In = log.New(os.Stdout, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info_Out = log.New(os.Stdout, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

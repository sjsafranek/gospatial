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
	serverLoggerWriter          io.Writer
	networkLoggerWriter   io.Writer
	network_logger_Info     *log.Logger
	network_logger_Warning  *log.Logger
	network_logger_Error    *log.Logger
	network_logger_Info_In  *log.Logger
	network_logger_Info_Out *log.Logger
)

func networkLoggerInit() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	log_file := strings.Replace(dir, "bin", "log/network.log", -1)
	networkLoggerWriter, err = os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	network_logger_Info_In = log.New(networkLoggerWriter, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info_Out = log.New(networkLoggerWriter, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info = log.New(networkLoggerWriter, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Warning = log.New(networkLoggerWriter, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Error = log.New(networkLoggerWriter, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func init() {
	AppMode = "logging"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	log_file := strings.Replace(dir, "bin", "log/server.log", -1)
	serverLoggerWriter, err := os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer serverLoggerWriter.Close()
	Info = log.New(serverLoggerWriter, "INFO  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(serverLoggerWriter, "DEBUG [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(serverLoggerWriter, "WARN  [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(serverLoggerWriter, "ERROR [SERVER] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)

	networkLoggerInit()
}

func testLoggerInit() {
	AppMode = "testing"
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
// 	network_logger_Info = log.New(os.Stdout, "INFO  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	network_logger_Warning = log.New(os.Stdout, "WARN  [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	network_logger_Error = log.New(os.Stderr, "ERROR [NETWORK] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	network_logger_Info_In = log.New(os.Stdout, "INFO  [NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// 	network_logger_Info_Out = log.New(os.Stdout, "INFO  [NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
// }

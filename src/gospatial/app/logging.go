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
	Trace                  *log.Logger
	Info                   *log.Logger
	Debug                  *log.Logger
	Warning                *log.Logger
	Error                  *log.Logger
	DebugModeLogFile       io.Writer
	network_logger_writer  io.Writer
	network_logger_Info    *log.Logger
	network_logger_Warning *log.Logger
	network_logger_Error   *log.Logger
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
	// network_logger_Info_In = log.New(network_logger_writer, "[NETWORK] [IN] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	// network_logger_Info_Out = log.New(network_logger_writer, "[NETWORK] [OUT] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info = log.New(network_logger_writer, "[NETWORK] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Warning = log.New(network_logger_writer, "[NETWORK] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Error = log.New(network_logger_writer, "[NETWORK] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func init() {
	// network_logger_init()
	Trace = log.New(os.Stdout, "[GOSPATIAL] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(os.Stdout, "[GOSPATIAL] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(os.Stdout, "[GOSPATIAL] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "[GOSPATIAL] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(os.Stderr, "[GOSPATIAL] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func DebugMode() {
	AppMode = "debug"
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	log_file := strings.Replace(dir, "bin", "log/server.log", -1)
	DebugModeLogFile, err := os.OpenFile(log_file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer DebugModeLogFile.Close()
	Trace = log.New(DebugModeLogFile, "[SERVER] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(DebugModeLogFile, "[SERVER] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(DebugModeLogFile, "[SERVER] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(DebugModeLogFile, "[SERVER] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(DebugModeLogFile, "[SERVER] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func test_logger_init() {
	AppMode = "test"
	DebugModeLogFile, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer DebugModeLogFile.Close()
	Trace = log.New(DebugModeLogFile, "[TESTING] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(DebugModeLogFile, "[TESTING] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(DebugModeLogFile, "[TESTING] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(DebugModeLogFile, "[TESTING] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(DebugModeLogFile, "[TESTING] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func StdOutMode() {
	AppMode = "standard"
	Trace = log.New(os.Stdout, "[SERVER] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(os.Stdout, "[SERVER] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(os.Stdout, "[SERVER] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "[SERVER] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(os.Stderr, "[SERVER] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Info = log.New(os.Stdout, "[NETWORK] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Warning = log.New(os.Stdout, "[NETWORK] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	network_logger_Error = log.New(os.Stderr, "[NETWORK] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace            *log.Logger
	Info             *log.Logger
	Debug            *log.Logger
	Warning          *log.Logger
	Error            *log.Logger
	DebugModeLogFile io.Writer
	WebClient        *log.Logger
)

func init() {
	Trace = log.New(os.Stdout, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(os.Stdout, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(os.Stdout, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(os.Stderr, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func DebugMode() {
	AppMode = "debug"
	DebugModeLogFile, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer DebugModeLogFile.Close()
	Trace = log.New(DebugModeLogFile, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(DebugModeLogFile, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(DebugModeLogFile, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(DebugModeLogFile, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(DebugModeLogFile, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func TestMode_old() {
	AppMode = "testing"
	Trace = log.New(ioutil.Discard, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(ioutil.Discard, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(ioutil.Discard, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(ioutil.Discard, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(ioutil.Discard, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func TestMode() {
	AppMode = "testing"
	DebugModeLogFile, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer DebugModeLogFile.Close()
	Trace = log.New(DebugModeLogFile, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(DebugModeLogFile, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(DebugModeLogFile, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(DebugModeLogFile, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(DebugModeLogFile, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func StandardMode() {
	AppMode = "standard"
	Trace = log.New(os.Stdout, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(os.Stdout, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Debug = log.New(os.Stdout, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(os.Stderr, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

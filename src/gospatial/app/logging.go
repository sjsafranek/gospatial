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
	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WebClient = log.New(os.Stdout, "CLIENT: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func DebugMode(status bool) {
	if status {
		AppMode = "debug"
		DebugModeLogFile, err := os.OpenFile("debug_mode.log", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			Error.Fatal("Error opening file: %v", err)
		}
		// defer DebugModeLogFile.Close()
		Trace = log.New(DebugModeLogFile, "TRACE : ", log.Ldate|log.Ltime|log.Lshortfile)
		Info = log.New(DebugModeLogFile, "INFO : ", log.Ldate|log.Ltime|log.Lshortfile)
		Debug = log.New(DebugModeLogFile, "DEBUG : ", log.Ldate|log.Ltime|log.Lshortfile)
		Warning = log.New(DebugModeLogFile, "WARNING : ", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(DebugModeLogFile, "ERROR : ", log.Ldate|log.Ltime|log.Lshortfile)
		WebClient = log.New(DebugModeLogFile, "CLIENT : ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {

		// Trace = log.New(os.Stdout, "TRACE : ", log.Ldate|log.Ltime|log.Lshortfile)
		// Info = log.New(os.Stdout, "INFO : ", log.Ldate|log.Ltime|log.Lshortfile)
		// Debug = log.New(os.Stdout, "DEBUG : ", log.Ldate|log.Ltime|log.Lshortfile)
		// Warning = log.New(os.Stdout, "WARNING : ", log.Ldate|log.Ltime|log.Lshortfile)
		// Error = log.New(os.Stderr, "ERROR : ", log.Ldate|log.Ltime|log.Lshortfile)
		// WebClient = log.New(os.Stdout, "CLIENT : ", log.Ldate|log.Ltime|log.Lshortfile)

		// Trace = log.New(os.Stdout, "[SERVER] TRACE : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		// Info = log.New(os.Stdout, "[SERVER] INFO : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		// Debug = log.New(os.Stdout, "[SERVER] DEBUG : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		// Warning = log.New(os.Stdout, "[SERVER] WARNING : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		// Error = log.New(os.Stderr, "[SERVER] ERROR : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		// WebClient = log.New(os.Stdout, "[SERVER] CLIENT : ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		Trace = log.New(os.Stdout, "[FIND] TRACE | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		Info = log.New(os.Stdout, "[FIND] INFO  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		Debug = log.New(os.Stdout, "[FIND] DEBUG | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		Warning = log.New(os.Stdout, "[FIND] WARN  | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, "[FIND] ERROR | ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		// WebClient = log.New(os.Stdout, "[FIND] CLIE |", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
		// CRITICAL
	}
}

package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var (
	InfoTcp         *log.Logger
	DebugTcp        *log.Logger
	WarningTcp      *log.Logger
	ErrorTcp        *log.Logger
	tcpLoggerWriter io.Writer
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	// server logging
	tcpLogFile := strings.Replace(dir, "bin", "log/socket.log", -1)
	tcpLoggerWriter, err := os.OpenFile(tcpLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Fatal("Error opening file: %v", err)
	}
	// defer tcpLoggerWriter.Close()
	InfoTcp = log.New(tcpLoggerWriter, "INFO  [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	DebugTcp = log.New(tcpLoggerWriter, "DEBUG [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	WarningTcp = log.New(tcpLoggerWriter, "WARN  [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	ErrorTcp = log.New(tcpLoggerWriter, "ERROR [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func TcpServer() {
	go func() {
		// Listen for incoming connections.
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			log.Println("Error listening:", err.Error())
			os.Exit(1)
		}
		// Close the listener when the application closes.
		defer l.Close()
		log.Println("Tcp Listening on " + CONN_HOST + ":" + CONN_PORT)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				ErrorTcp.Println("Error accepting: ", err.Error())
				// os.Exit(1)
				return
			}
			InfoTcp.Println("Connection open")
			// Handle connections in a new goroutine.
			go handleTcpClient(conn)
		}
	}()
}

// Handles incoming requests.
func handleTcpClient(conn net.Conn) {
	
	defer conn.Close()
	
	for {
		
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		
		// output message received
		InfoTcp.Print("Message Received:", string(message))
		
		// json parse message
		req := make(map[string]interface{})
		err := json.Unmarshal([]byte(message), &req)
		if err != nil {
			// invalid message
			// close connection
			WarningTcp.Println("error:", err)
			resp := `{"status": "error", "error": "` + fmt.Sprintf("%v", err) + `"}`
			conn.Write([]byte(resp + "\n"))
			InfoTcp.Println("Connection closed")
			return
		}

		// get method
		success := false
		switch {
		    case req["method"] == "create_customer":
				resp := `{"status": "success", "data": {}}`
				conn.Write([]byte(resp + "\n"))
				success = true
		    case req["method"] == "create_layer":
				resp := `{"status": "success", "data": {}}`
				conn.Write([]byte(resp + "\n"))
				success = true
		    case req["method"] == "create_feature":
				resp := `{"status": "success", "data": {}}`
				conn.Write([]byte(resp + "\n"))
				success = true
	    }
	    if !success {
	    	resp := `{"status": "error", "error": "method not found"}`
	    	conn.Write([]byte(resp + "\n"))
	    }

	}
}

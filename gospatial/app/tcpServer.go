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
	TCP_DEFAULT_CONN_HOST = "localhost"
	TCP_DEFAULT_CONN_PORT = "3333"
	// TCP_DEFAULT_CONN_TYPE = "tcp"
)

var (
	infoTcp         *log.Logger
	debugTcp        *log.Logger
	warningTcp      *log.Logger
	errorTcp        *log.Logger
	tcpLoggerWriter io.Writer
)


type TcpServer struct {
	Host string
	Port string
}

func (self TcpServer) init() {
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
	infoTcp = log.New(tcpLoggerWriter, "INFO  [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	debugTcp = log.New(tcpLoggerWriter, "DEBUG [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	warningTcp = log.New(tcpLoggerWriter, "WARN  [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	errorTcp = log.New(tcpLoggerWriter, "ERROR [TCP] ", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
}

func (self TcpServer) Start() {
	self.init()
	go func() {
		// Check settings and apply defaults
		host := self.Host
		if host == "" {
			host = TCP_DEFAULT_CONN_HOST
		}

		port := self.Port
		if port == "" {
			port = TCP_DEFAULT_CONN_PORT
		}

		// Listen for incoming connections.
		l, err := net.Listen("tcp", host+":"+port)
		if err != nil {
			log.Println("Error listening:", err.Error())
			errorTcp.Println("Error listening:", err.Error())
			panic(err)
		}

		// Close the listener when the application closes.
		defer l.Close()
		log.Println("Tcp Listening on " + host + ":" + port)
		// infoTcp.Println("Error listening:", err.Error())
		for {

			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				errorTcp.Println("Error accepting: ", err.Error())
				return
			}

			infoTcp.Println("Connection open", conn.RemoteAddr().String())

			// check for local connection
			if strings.Contains(conn.RemoteAddr().String(), "127.0.0.1") {
				// Handle connections in a new goroutine.
				go self.tcpClientHandler(conn)
			} else {
				conn.Close()
			}

		}
	}()
}

// Handles incoming requests.
func (self TcpServer) tcpClientHandler(conn net.Conn) {

	defer conn.Close()

	for {

		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')

		// output message received
		infoTcp.Println("Message Received: ", string(message))

		// json parse message
		req := make(map[string]interface{})
		err := json.Unmarshal([]byte(message), &req)
		if err != nil {
			// invalid message
			// close connection
			warningTcp.Println("error:", err)
			resp := `{"status": "error", "error": "` + fmt.Sprintf("%v", err) + `"}`
			conn.Write([]byte(resp + "\n"))
			infoTcp.Println("Connection closed")
			return
		}

		// get method
		success := false
		switch {
		case req["method"] == "clear_datasource_cache":
			// {"method": "clear_datasource_cache"}
			// Unload all layers in database cache
			for key := range DB.Cache {
				delete(DB.Cache, key)
			}
			resp := `{"status": "success", "data": {}}`
			conn.Write([]byte(resp + "\n"))
			success = true
		case req["method"] == "loaded_datasources":
			// {"method": "loaded_datasources"}
			// result := make(map[string]interface{})
			result, _ := json.Marshal(DB.Cache)
			resp := `{"status": "success", "data": ` + string(result) + `}`
			conn.Write([]byte(resp + "\n"))
			success = true
		case req["method"] == "clear_customer_cache":
			// {"method": "clear_customer_cache"}
			// Unload all apikeys in database cache
			for key := range DB.Apikeys {
				delete(DB.Apikeys, key)
			}
			resp := `{"status": "success", "data": {}}`
			conn.Write([]byte(resp + "\n"))
			success = true
		case req["method"] == "share_datasource":
			// req["data"]["datasource"]
			// req["data"]["to"]
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

package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"gospatial/utils"
)

const (
	TCP_DEFAULT_CONN_HOST = "localhost"
	TCP_DEFAULT_CONN_PORT = "3333"
	// TCP_DEFAULT_CONN_TYPE = "tcp"
)


type TcpServer struct {
	Host string
	Port string
}

func (self TcpServer) Start() {
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
			ServerLogger.Error("Error listening:", err.Error())
			panic(err)
		}

		// Close the listener when the application closes.
		defer l.Close()
		// log.Println("Tcp Listening on " + host + ":" + port)
		ServerLogger.Info("Tcp Listening on " + host + ":" + port)

		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				NetworkLogger.Error("Error accepting: ", err.Error())
				return
			}

			NetworkLogger.Info("Connection open ", conn.RemoteAddr().String(), " [TCP]")

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

	authenticated := false

	for {

		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')

		// output message received
		NetworkLogger.Info("Message Received: ", string(message), " [TCP]")

		// json parse message
		// req := make(map[string]string)
		req := TcpMessage{}
		err := json.Unmarshal([]byte(message), &req)
		if err != nil {
			// invalid message
			// close connection
			NetworkLogger.Warn("error:", err)
			resp := `{"status": "error", "error": "` + fmt.Sprintf("%v", err) + `"}`
			conn.Write([]byte(resp + "\n"))
			NetworkLogger.Info("Connection closed"," [TCP]")
			return
		}

		// get method
		if !authenticated {
			if req.Method == "authenticate" {
				// {"method":"authenticate", "authkey": "O1p9dLhsryIn"}
				authenticated = SuperuserKey == req.Authkey
				if authenticated {
					resp := `{"status": "success", "data": {}}`
					conn.Write([]byte(resp + "\n"))
				} else {
					NetworkLogger.Warn("error: incorrect authkey", " [TCP]")
					resp := `{"status": "error", "error": "incorrect authkey"}`
					conn.Write([]byte(resp + "\n"))
				}
			} else {
				resp := `{"status": "error", "error": "connection not authenticated"}`
				conn.Write([]byte(resp + "\n"))
			}
		} else {

			success := false
			switch {
			case req.Method == "clear_datasource_cache" && authenticated:
				// {"method": "clear_datasource_cache"}
				// Unload all layers in database cache
				for key := range DB.Cache {
					delete(DB.Cache, key)
				}
				resp := `{"status": "success", "data": {}}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "loaded_datasources" && authenticated:
				// {"method": "loaded_datasources"}
				result, _ := json.Marshal(DB.Cache)
				resp := `{"status": "success", "data": ` + string(result) + `}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "clear_customer_cache" && authenticated:
				// {"method": "clear_customer_cache"}
				// Unload all apikeys in database cache
				for key := range DB.Apikeys {
					delete(DB.Apikeys, key)
				}
				resp := `{"status": "success", "data": {}}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "assign_datasource" && authenticated:
				datasource_id := req.Datasource //["datasource_id"]
				apikey := req.Apikey //["apikey"]
				customer, err := DB.GetCustomer(apikey)
				resp := `{"status": "success", "data": {}}`
				if err != nil {
					fmt.Println("Customer key not found!")
					resp = `{"status": "error", "data": {"error": "` + err.Error() + `", "message": "Customer key not found!"}}`
				}
				// CHECK IF DATASOURCE EXISTS
				// *****
				fmt.Println(DB.GetLayer(datasource_id))

				customer.Datasources = append(customer.Datasources, datasource_id)
				DB.InsertCustomer(customer)
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "create_user" && authenticated:
				apikey := utils.NewAPIKey(12)
				customer := Customer{Apikey: apikey}
				resp := `{"status": "success", "data": {"apikey": "` + apikey + `"}}`
				err := DB.InsertCustomer(customer)
				if err != nil {
					fmt.Println(err)
					resp = `{"status": "error", "data": {"error": "` + err.Error() + `", "message": "error creating customer"}}`
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			//  Replay database
			case req.Method == "insert_apikey" && authenticated:
				customer := Customer{Apikey: req.Data.Apikey, Datasources: req.Data.Datasources}
				resp := `{"status": "success", "data": {"apikey": "` + req.Data.Apikey + `"}}`
				err := DB.InsertCustomer(customer)
				if err != nil {
					fmt.Println(err)
					resp = `{"status": "error", "data": {"error": "` + err.Error() + `", "message": "error creating customer"}}`
				}
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "insert_feature" && authenticated:
				err = DB.InsertFeature(req.Data.Datasource, req.Data.Feature)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","datasource":"` + req.Data.Datasource + `", "message":"feature added"}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "new_layer" && authenticated:
				err = DB.InsertLayer(req.Data.Datasource, req.Data.Layer)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","datasource":"` + req.Data.Datasource + `"}`
				conn.Write([]byte(resp + "\n"))
				success = true

		/*
			case req.Metho == "delete_layer" && authenticated:
				req.Data.Datasource
		*/

			}

			if !success {
				resp := `{"status": "error", "error": "method not found"}`
				conn.Write([]byte(resp + "\n"))
			}

		}

	}
}

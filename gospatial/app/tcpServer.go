package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gospatial/utils"
	"net"
	"net/textproto"
	"strings"
)

const (
	TCP_DEFAULT_CONN_HOST = "localhost"
	TCP_DEFAULT_CONN_PORT = "3333"
	TCP_DEFAULT_CONN_TYPE = "tcp"
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
		l, err := net.Listen(TCP_DEFAULT_CONN_TYPE, host+":"+port)
		if err != nil {
			ServerLogger.Error("Error listening:", err.Error())
			panic(err)
		}

		// Close the listener when the application closes.
		defer l.Close()

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

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	defer conn.Close()

	// DEBUGGING
	//	authenticated := false
	authenticated := true

	for {

		// will listen for message to process ending in newline (\n)
		//message, _ := bufio.NewReader(conn).ReadString('\n') // sometimes read partial messages
		message, _ := tp.ReadLine()

		// go self.handleTcpMessage(authenticated, message, conn)

		// output message received
		NetworkLogger.Info("Message Received: ", string(message), " [TCP]")

		// json parse message
		req := TcpMessage{}
		err := json.Unmarshal([]byte(message), &req)
		if err != nil {
			// invalid message
			// close connection
			NetworkLogger.Warn("error:", err)
			resp := `{"status": "error", "error": "` + fmt.Sprintf("%v", err) + `",""}`
			conn.Write([]byte(resp + "\n"))
			NetworkLogger.Info("Connection closed", " [TCP]")
			return
		}

		// get method
		if !authenticated {
			if req.Method == "authenticate" {
				// {"method":"authenticate", "authkey": "7q1qcqmsxnvw"}
				authenticated = SuperuserKey == req.Authkey
				if authenticated {
					resp := `{"status": "ok", "data": {}}`
					conn.Write([]byte(resp + "\n"))
				} else {
					NetworkLogger.Warn("error: incorrect authkey", " [TCP]")
					resp := `{"status": "error", "error": "incorrect authkey"}`
					conn.Write([]byte(resp + "\n"))
				}
			} else if req.Method == "help" {
				conn.Write([]byte("Methods:\n"))
				conn.Write([]byte("\t authenticate\n"))
				conn.Write([]byte("\t create_user\n"))
				conn.Write([]byte("\t insert_apikey\n"))
				conn.Write([]byte("\t export_apikeys\n"))
				conn.Write([]byte("\t export_apikey\n"))
				conn.Write([]byte("\t new_layer\n"))
				conn.Write([]byte("\t insert_feature\n"))
				conn.Write([]byte("\t export_datasource\n"))
				conn.Write([]byte("\t export_datasources\n"))
			} else {
				resp := `{"status": "error", "error": "connection not authenticated"}`
				conn.Write([]byte(resp + "\n"))
			}
		} else {

			success := false
			switch {

			case req.Method == "assign_datasource" && authenticated:
				datasource_id := req.Datasource //["datasource_id"]
				apikey := req.Apikey            //["apikey"]
				customer, err := DB.GetCustomer(apikey)
				resp := `{"status": "ok", "data": {}}`
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

			case req.Method == "create_apikey" && authenticated:
				// {"method":"create_user"}
				apikey := utils.NewAPIKey(12)
				customer := Customer{Apikey: apikey}
				resp := `{"status": "ok", "data": {"apikey": "` + apikey + `"}}`
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
				resp := `{"status": "ok", "data": {"apikey": "` + req.Data.Apikey + `"}}`
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

			case req.Method == "edit_feature" && authenticated:
				err = DB.EditFeature(req.Data.Datasource, req.Data.GeoId, req.Data.Feature)
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

			case req.Method == "export_apikeys" && authenticated:
				// {"method":"export_apikeys"}
				apikeys, err := DB.SelectAll("apikeys")
				if err != nil {
					fmt.Println(err)
				}
				js, err := json.Marshal(apikeys)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","data":"` + string(js) + `"}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_apikey" && authenticated:
				// {"method":"export_apikey","apikey":"4AvJJ3oW0zeT"}
				apikey, err := DB.GetCustomer(req.Apikey)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(apikey)
				js, err := json.Marshal(apikey)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","data":"` + string(js) + `"}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_datasources" && authenticated:
				// {"method":"export_datasources"}
				layers, err := DB.SelectAll("layers")
				if err != nil {
					fmt.Println(err)
				}
				js, err := json.Marshal(layers)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","data":"` + string(js) + `"}`
				conn.Write([]byte(resp + "\n"))
				success = true

			case req.Method == "export_datasource" && authenticated:
				// {"method":"export_datasource","datasource":"3b1f5d633d884b9499adfc9b49c45236"}
				layer, err := DB.GetLayer(req.Datasource)
				if err != nil {
					fmt.Println(err)
				}
				js, err := json.Marshal(layer)
				if err != nil {
					fmt.Println(err)
				}
				resp := `{"status":"ok","data":"` + string(js) + `"}`
				conn.Write([]byte(resp + "\n"))
				success = true

			}

			if !success {
				resp := `{"status": "error", "error": "method not found"}`
				conn.Write([]byte(resp + "\n"))
			}

		}

	}
}

/*

"insert_layer"
"delete_layer"

*/

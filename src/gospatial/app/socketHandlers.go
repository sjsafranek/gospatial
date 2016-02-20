package app

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	ws *websocket.Conn
	ds string
}

type hub struct {
	Sockets map[*websocket.Conn]string
}

func (h hub) broadcast(conn *connection) {
	Trace.Println("Broadcasting message to open websocket connections")
	for i := range h.Sockets {
		if i != conn.ws && h.Sockets[i] == conn.ds {
			i.WriteMessage(websocket.TextMessage, []byte(`update layer`))
		}
	}
}

var Hub = hub{
	Sockets: make(map[*websocket.Conn]string),
}

func messageListener(conn *connection) {
	defer func() {
		Debug.Printf("Disconnecting websocket connection")
		conn.ws.Close()
		delete(Hub.Sockets, conn.ws)
	}()
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			Warning.Println(err)
			break
		}
		Debug.Printf("Message: %s %s", string(message), conn.ds)
		Hub.broadcast(conn)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	Trace.Println("Establishing websocket connection")
	vars := mux.Vars(r)
	ds := vars["ds"]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Info.Println(err)
		return
	}
	conn := connection{ws: ws, ds: ds}
	Hub.Sockets[ws] = ds
	go messageListener(&conn)
}

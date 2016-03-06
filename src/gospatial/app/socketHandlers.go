package app

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	ws *websocket.Conn
	ds string
	ip string
	c  int
}

type hub struct {
	Sockets map[string]map[int]*websocket.Conn
}

func (h hub) broadcast(update bool, conn *connection) {
	type Message struct {
		Update  bool `json:"update"`
		Viewers int  `json:"viewers"`
	}
	Trace.Println("WebSockets | Broadcasting message to open connections")
	msg := Message{Update: update, Viewers: len(h.Sockets[conn.ds])}
	for i := range h.Sockets[conn.ds] {
		if h.Sockets[conn.ds][i] != conn.ws {
			Trace.Println("WebSockets | Sending message to client")
			h.Sockets[conn.ds][i].WriteJSON(msg)
		}
	}
}

func (h hub) broadcastAllDsViewers(update bool, ds string) {
	type Message struct {
		Update  bool `json:"update"`
		Viewers int  `json:"viewers"`
	}
	Trace.Println("WebSockets | Broadcasting message to open connections")
	msg := Message{Update: update, Viewers: len(h.Sockets[ds])}
	for i := range h.Sockets[ds] {
		Trace.Println("WebSockets | Sending message to client")
		h.Sockets[ds][i].WriteJSON(msg)
	}
}

var Hub = hub{
	Sockets: make(map[string]map[int]*websocket.Conn),
}

func messageListener(conn *connection) {
	defer func() {
		Debug.Printf("WebSockets | Disconnecting %s connection", conn.ip)
		conn.ws.Close()
		delete(Hub.Sockets[conn.ds], conn.c)
		Hub.broadcast(false, conn)
	}()
	for {
		var m interface{}
		err := conn.ws.ReadJSON(&m)
		if err != nil {
			Error.Printf("%s | %s", conn.ip, err)
			return
		}
		Debug.Printf("WebSockets | Message: %v %s", m, conn.ds)
		for i := range Hub.Sockets[conn.ds] {
			if Hub.Sockets[conn.ds][i] != conn.ws {
				Trace.Println("WebSockets | Sending message to client")
				Hub.Sockets[conn.ds][i].WriteJSON(m)
			}
		}
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
	Trace.Println(r.RemoteAddr, "| Establishing websocket connection")
	vars := mux.Vars(r)
	ds := vars["ds"]
	ip := r.RemoteAddr
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Error.Println(err)
		return
	}
	conn := connection{ws: ws, ds: ds, ip: ip, c: len(Hub.Sockets[ds])}
	if _, ok := Hub.Sockets[ds]; ok {
		Hub.Sockets[ds][len(Hub.Sockets[ds])] = ws
	} else {
		Hub.Sockets[ds] = make(map[int]*websocket.Conn)
		Hub.Sockets[ds][conn.c] = ws
		Info.Println(r.RemoteAddr, "| WebSocket connection open")
	}
	Hub.broadcastAllDsViewers(false, conn.ds)
	go messageListener(&conn)
}

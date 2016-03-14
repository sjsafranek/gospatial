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

// Websocket status codes
// http://tools.ietf.org/html/rfc6455#page-45

func (self hub) broadcast(update bool, conn *connection) {
	type Message struct {
		Update  bool `json:"update"`
		Viewers int  `json:"viewers"`
	}
	Trace.Println("Broadcasting message to open connections")
	msg := Message{Update: update, Viewers: len(self.Sockets[conn.ds])}
	for i := range self.Sockets[conn.ds] {
		if self.Sockets[conn.ds][i] != conn.ws {
			Trace.Println("Sending message to client")
			self.Sockets[conn.ds][i].WriteJSON(msg)
		}
	}
}

func (self hub) broadcastAllDsViewers(update bool, ds string) {
	type Message struct {
		Update  bool `json:"update"`
		Viewers int  `json:"viewers"`
	}
	Trace.Println("Broadcasting message to open connections")
	msg := Message{Update: update, Viewers: len(self.Sockets[ds])}
	for i := range self.Sockets[ds] {
		Trace.Println("Sending message to client")
		self.Sockets[ds][i].WriteJSON(msg)
	}
}

var Hub = hub{
	Sockets: make(map[string]map[int]*websocket.Conn),
}

func messageListener(conn *connection) {
	defer func() {
		conn.ws.Close()
		delete(Hub.Sockets[conn.ds], conn.c)
		if len(Hub.Sockets[conn.ds]) == 0 {
			delete(Hub.Sockets, conn.ds)
		}
		Hub.broadcast(false, conn)
	}()
	for {
		var m interface{}
		err := conn.ws.ReadJSON(&m)
		if err != nil {
			Error.Println(conn.ip, "WS /ws/"+conn.ds+" [1001]")
			Error.Printf("%s %s", conn.ip, err)
			return
		}
		Debug.Printf("Message: %v %s", m, conn.ds)
		for i := range Hub.Sockets[conn.ds] {
			if Hub.Sockets[conn.ds][i] != conn.ws {
				Trace.Println("Sending message to client")
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
	vars := mux.Vars(r)
	ds := vars["ds"]
	ip := r.RemoteAddr
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Error.Println(r.RemoteAddr, "WS /ws/"+ds+" [500]")
		Error.Println(err)
		return
	}
	conn := connection{ws: ws, ds: ds, ip: ip, c: len(Hub.Sockets[ds])}
	if _, ok := Hub.Sockets[ds]; ok {
		Hub.Sockets[ds][len(Hub.Sockets[ds])] = ws
	} else {
		Hub.Sockets[ds] = make(map[int]*websocket.Conn)
		Hub.Sockets[ds][conn.c] = ws
		// Info.Println(r.RemoteAddr, "WS Open [200]")
		Info.Println(r.RemoteAddr, "WS /ws/"+conn.ds+" [200]")
	}
	Hub.broadcastAllDsViewers(false, conn.ds)
	go messageListener(&conn)
}

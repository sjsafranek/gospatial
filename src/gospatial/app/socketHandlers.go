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

func (h hub) broadcast(c *connection) {
	for i := range h.Sockets {
		Info.Println(h.Sockets[i])
		if i != c.ws && h.Sockets[i] == c.ds {
			i.WriteMessage(websocket.TextMessage, []byte(`message`))
		}
	}
}

var Hub = hub{
	Sockets: make(map[*websocket.Conn]string),
}

func messageListener(conn *connection) {
	for {
		defer func() {
			conn.ws.Close()
			delete(Hub.Sockets, conn.ws)
		}()
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			Error.Println(err)
			break
		}
		Info.Println(string(message))
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

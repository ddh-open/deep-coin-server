package service

import (
	"devops-http/app/module/ws/model"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWs(conn *websocket.Conn) *websocketService {
	return &websocketService{conn: conn, m: sync.Mutex{}}
}

type websocketService struct {
	conn   *websocket.Conn
	m      sync.Mutex
	closed bool
}

func (ws *websocketService) ReadMessage() (int, []byte, error) {
	ws.m.Lock()
	defer ws.m.Unlock()
	return ws.conn.ReadMessage()
}

func (ws *websocketService) ReadMessageUntilHas() (int, []byte, error) {
	ws.m.Lock()
	defer ws.m.Unlock()
	//for {
	//	_, data, err := ws.conn.ReadMessage()
	//}
	return ws.conn.ReadMessage()
}

func (ws *websocketService) WriteMessage(messageType int, message []byte) error {
	ws.m.Lock()
	defer ws.m.Unlock()
	return ws.conn.WriteMessage(messageType, message)
}

func (ws *websocketService) Close() {
	ws.m.Lock()
	ws.closed = true
	ws.conn.Close()
	ws.m.Unlock()
}

func HandleMessageError(ws *websocketService, err error) {
	ws.WriteMessage(1, []byte(err.Error()))
}

func HandleMessage(ws *websocketService, data model.WebSocketReadMessage) {
	switch data.Type {
	case "host-shell":

	}
}

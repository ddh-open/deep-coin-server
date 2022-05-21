package def

import (
	"devops-http/app/module/ws/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
)

func NewWs(conn *websocket.Conn) *WebsocketService {
	return &WebsocketService{conn: conn, m: &sync.Mutex{}}
}

type WebsocketService struct {
	conn   *websocket.Conn
	m      *sync.Mutex
	closed bool
}

func (ws *WebsocketService) ReadMessage() (messageType int, data []byte, err error) {
	messageType, data, err = ws.conn.ReadMessage()
	return
}

func (ws *WebsocketService) WriteMessage(messageType int, message []byte) error {
	ws.m.Lock()
	err := ws.conn.WriteMessage(messageType, message)
	ws.m.Unlock()
	return err
}

func (ws *WebsocketService) Close() {
	ws.m.Lock()
	ws.closed = true
	ws.conn.Close()
	ws.m.Unlock()
}

func (ws *WebsocketService) HandleMessageError(err error, uuid *uuid.UUID) {
	errData := model.WebSocketWriteMessage{
		UUID: uuid,
		Type: "error",
		Data: err.Error(),
	}
	errByteData, _ := json.Marshal(&errData)
	err = ws.WriteMessage(1, errByteData)
}

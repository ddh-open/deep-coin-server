package service

import (
	"devops-http/app/module/ws/model"
	"devops-http/app/module/ws/service/def"
	"devops-http/app/module/ws/service/hostShell"
	"github.com/gorilla/websocket"
	"net/http"
)

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleMessage(ws *def.WebsocketService, data model.WebSocketReadMessage) {
	var err error
	switch data.Type {
	case "host-shell":
		err = hostShell.WsHostShellHandle(ws, &data)
	}

	if err != nil {
		ws.HandleMessageError(err, data.UUID)
	}
}

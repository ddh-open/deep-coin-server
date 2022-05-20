package model

import uuid "github.com/satori/go.uuid"

// WebSocketReadMessage websocket 相关
type WebSocketReadMessage struct {
	UUID   *uuid.UUID        `json:"uuid"`
	Type   string            `json:"type"`
	Action string            `json:"action"`
	Param  map[string]string `json:"param"`
}

type WebSocketWriteMessage struct {
	UUID *uuid.UUID `json:"uuid"`
	Type string     `json:"type"`
	Data string     `json:"data"`
}

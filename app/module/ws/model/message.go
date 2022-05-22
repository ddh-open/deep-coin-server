package model

// WebSocketReadMessage websocket 相关
type WebSocketReadMessage struct {
	UUID   string            `json:"uuid"`
	Type   string            `json:"type"`
	Action string            `json:"action"`
	Param  map[string]string `json:"param"`
}

type WebSocketWriteMessage struct {
	UUID string `json:"uuid"`
	Type string `json:"type"`
	Data string `json:"data"`
}

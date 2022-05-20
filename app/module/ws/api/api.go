package api

import (
	"devops-http/app/module/base/response"
	"devops-http/app/module/ws/model"
	"devops-http/app/module/ws/service"
	"devops-http/framework/gin"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

func Register(r *gin.Engine) error {
	r.GET("base/ws", CreateWs)
	return nil
}

// CreateWs godoc
// @Summary 获得全局通用ws
// @Security ApiKeyAuth
// @Description 获得全局通用ws
// @accept application/json
// @Produce application/json
// @Tags Ws
// @Success 200 {object}  response.Response
// @Router /sys/ws [get]
func CreateWs(c *gin.Context) {
	res := response.Response{Code: 1, Msg: "", Data: nil}
	// 升级get请求为webSocket协议
	conn, err := service.UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		return
	}
	ws := service.NewWs(conn)
	defer ws.Close()
	var readMessage model.WebSocketReadMessage
	for {
		_, data, err := ws.ReadMessage()
		if err == nil {
			err = json.Unmarshal(data, &readMessage)
			if err != nil {
				go service.HandleMessageError(ws, err)
				continue
			}
			// 处理真实的业务功能
			if readMessage.UUID == nil {
				// 第一次来分配个uuid
				*(readMessage.UUID) = uuid.NewV1()
			}
			// 把业务功能分配出去
			go service.HandleMessage(ws, readMessage)
		}
	}
}

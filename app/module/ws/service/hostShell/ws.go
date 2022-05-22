package hostShell

import (
	"context"
	"devops-http/app/module/ws/model"
	"devops-http/app/module/ws/service/def"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func WsHostShellHandle(ws *def.WebsocketService, data *model.WebSocketReadMessage) (err error) {
	links := NewLinkManage()
	shell := links.GetLink(data.UUID)
	var sendData model.WebSocketWriteMessage
	sendData.Type = data.Type
	sendData.UUID = data.UUID
	// 如果不存在， 就先创建连接
	if shell == nil {
		if data.Param != nil && data.Param["user"] != "" && data.Param["ip"] != "" && data.Param["port"] != "" {
			shell = NewContext(data.Param["ip"], cast.ToInt(data.Param["port"]), data.Param["user"])
			err = shell.InitTerminalWithPassword("dou.190824")
			if err != nil {
				return
			}
			links.AddLink(data.UUID, shell)
		} else {
			err = errors.New("未连接web-shell, 参数错误")
			return err
		}
	}

	if shell == nil {
		err = errors.New("动作action create shell终端失败")
	}

	// 执行脚本
	if data.Param["shell"] != "" {
		ctx, cancel := context.WithCancel(context.Background())
		// 读执行的脚本日志
		go func() {
		outLabel:
			for {
				select {
				case <-ctx.Done():
					// 已经结束
					break outLabel
				case msg := <-shell.Logs:
					sendData.Data = msg
					sendByteData, _ := json.Marshal(&sendData)
					ws.WriteMessage(1, sendByteData)
				}
			}
		}()
		err = shell.SendCmd(data.Param["shell"])
		cancel()
	} else {
	outLabel:
		for {
			select {
			case msg := <-shell.Logs:
				sendData.Data = msg
				sendByteData, _ := json.Marshal(&sendData)
				ws.WriteMessage(1, sendByteData)
			default:
				if shell.Start {
					break outLabel
				}
			}
		}
	}
	return err
}

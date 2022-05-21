package ws

import (
	"devops-http/app/module/ws/api"
	"devops-http/framework/gin"
)

func Register(r *gin.Engine) error {
	r.GET("base/ws", api.CreateWs)
	return nil
}

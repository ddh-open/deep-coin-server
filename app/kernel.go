package app

import (
	"devops-http/boot"
	"devops-http/framework"
	"devops-http/framework/contract"
	"devops-http/framework/gin"
	"devops-http/framework/middleware"
)

// NewHttpEngine 创建了一个绑定了路由的Web引擎
func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	// 设置为Release，为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个Web引擎
	r := gin.New()
	// 设置了Engine
	r.SetContainer(container)
	// 默认注册recovery中间件
	r.Use(middleware.RecoveryWithZap(), middleware.GinLogger())
	// 生成swagger
	if container.MustMake(contract.EnvKey).(contract.Env).Get("APP_ENV") == "dev" {
		boot.SwaggerInit()
	}
	// 业务绑定路由操作
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}

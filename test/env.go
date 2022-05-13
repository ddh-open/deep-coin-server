package test

import (
	"devops-http/framework"
	"devops-http/framework/provider/app"
	"devops-http/framework/provider/config"
	"devops-http/framework/provider/env"
)

const (
	BasePath = "/Users/freemud/Desktop/devops-grpc/"
)

func InitBaseContainer() *framework.NiceContainer {
	// 初始化服务容器
	container := framework.NewNiceContainer()
	// 绑定App服务提供者
	container.Bind(&app.NiceAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.NiceEnvProvider{})
	container.Bind(&config.NiceConfigProvider{})
	return container
}

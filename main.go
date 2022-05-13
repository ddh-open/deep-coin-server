package main

import (
	app2 "devops-http/app"
	"devops-http/app/provider/casbin"
	"devops-http/app/provider/grpc"
	"devops-http/app/provider/jwt"
	"devops-http/app/provider/ldap"
	"devops-http/boot"
	"devops-http/framework"
	"devops-http/framework/provider/app"
	"devops-http/framework/provider/cache"
	"devops-http/framework/provider/config"
	"devops-http/framework/provider/env"
	"devops-http/framework/provider/id"
	"devops-http/framework/provider/kernel"
	"devops-http/framework/provider/log"
	"devops-http/framework/provider/orm"
	"devops-http/framework/provider/redis"
	"devops-http/framework/provider/trace"
)

// @title Swagger Devops API
// @version 0.0.1
// @description This is Devops API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name bearer
// @BasePath /
func main() {
	// 初始化服务容器
	container := framework.NewNiceContainer()
	// 绑定App服务提供者
	container.Bind(&app.NiceAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.NiceEnvProvider{})
	container.Bind(&config.NiceConfigProvider{})
	container.Bind(&id.NiceIDProvider{})
	container.Bind(&trace.NiceTraceProvider{})
	container.Bind(&log.NiceLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&redis.ProviderRedis{})
	container.Bind(&cache.NiceCacheProvider{})
	// 业务服务
	container.Bind(&grpc.ProviderGrpc{})
	container.Bind(&casbin.ProviderCabin{})
	container.Bind(&ldap.ProviderLdap{})
	container.Bind(&jwt.ProviderJWT{})
	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := app2.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.NiceKernelProvider{HttpEngine: engine})
	}
	boot.InitService(container)
}

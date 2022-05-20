package app

import (
	"devops-http/app/middleware"
	"devops-http/app/module/cluster/api/host"
	"devops-http/app/module/sys/api/config"
	"devops-http/app/module/sys/api/domain"
	"devops-http/app/module/sys/api/group"
	"devops-http/app/module/sys/api/icon"
	"devops-http/app/module/sys/api/menu"
	"devops-http/app/module/sys/api/operation"
	"devops-http/app/module/sys/api/path"
	"devops-http/app/module/sys/api/role"
	"devops-http/app/module/sys/api/user"
	"devops-http/app/module/third/api/apm"
	"devops-http/app/module/third/api/cls"
	"devops-http/app/module/third/api/tencent"
	workflow "devops-http/app/module/workflow/api"
	"devops-http/app/swagger"
	"devops-http/framework/contract"
	"devops-http/framework/gin"
	ginSwagger "devops-http/framework/middleware/gin-swagger"
	"devops-http/framework/middleware/gin-swagger/swaggerFiles"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	// set swagger info
	swagger.SwaggerInfo.Title = "Devops-Http API"
	swagger.SwaggerInfo.Description = "This is Devops-Http API"
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = ""
	swagger.SwaggerInfo.BasePath = ""
	swagger.SwaggerInfo.Schemes = []string{"http"}
	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.Auth(), middleware.OperationRecord())
	/** 系统相关  **/
	domain.Register(r)
	group.Register(r)
	menu.Register(r)
	path.Register(r)
	role.Register(r)
	operation.Register(r)
	config.Register(r)
	icon.Register(r)
	// 用户模块注册路由
	user.Register(r)
	// cmdb 主机模块
	host.Register(r)

	/** 第三方相关  **/
	apm.Register(r)
	cls.Register(r)
	tencent.Register(r)

	/** 工作流 **/
	workflow.Register(r)
}

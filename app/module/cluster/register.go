package cluster

import (
	"devops-http/app/module/cluster/api/application"
	"devops-http/app/module/cluster/api/host"
	"devops-http/framework/gin"
)

func Register(r *gin.Engine) error {
	apiApp := application.NewApplicationApi(r.GetContainer())
	appGroup := r.Group("/cluster/", func(c *gin.Context) {
	})
	appGroup.POST("application/list", apiApp.List)
	appGroup.GET("application/:id", apiApp.GetApplication)
	appGroup.POST("application/save", apiApp.AddApplication)
	appGroup.PUT("application/modify", apiApp.Modify)
	appGroup.DELETE("application/delete/:id", apiApp.Delete)

	appGroup.POST("application/add/config", apiApp.AddApplicationConfig)
	appGroup.PUT("application/modify/config", apiApp.ModifyConfig)
	appGroup.DELETE("application/delete/config", apiApp.DeleteConfig)

	apiHost := host.NewHostApi(r.GetContainer())
	cmdbGroup := r.Group("/cmdb/", func(c *gin.Context) {
	})
	cmdbGroup.POST("host/list", apiHost.GetHostList)
	cmdbGroup.GET("host/group/tree", apiHost.GetHostGroupTree)
	return nil
}

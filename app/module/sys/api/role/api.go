package role

import (
	"devops-http/app/module/sys/service/role"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/gin"
)

type ApiRole struct {
	service *role.Service
}

func Register(r *gin.Engine) error {
	api := NewSysApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})

	// 用户角色相关接口
	sysGroup.GET("roles/:id", api.GetRoles)
	sysGroup.POST("roles/list", api.ListRoles)
	sysGroup.POST("roles/add", api.AddRole)
	sysGroup.POST("roles/add/resources", api.AddResourcesToRole)
	sysGroup.POST("roles/modify", api.ModifyRole)
	sysGroup.POST("roles/copy", api.CopyRole)
	sysGroup.DELETE("roles/delete", api.DeleteRole)

	return nil
}

func NewSysApi(c framework.Container) *ApiRole {
	return &ApiRole{service: role.NewService(c)}
}

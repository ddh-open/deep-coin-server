package group

import (
	"devops-http/app/module/sys/service/group"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiGroup struct {
	service *group.Service
}

func Register(r *gin.Engine) error {
	api := NewGroupApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})

	// 用户组相关接口
	sysGroup.GET("group/:id", api.GetGroups)
	sysGroup.POST("group/list", api.ListGroups)
	sysGroup.GET("group/tree", api.TreeGroups)
	sysGroup.POST("group/add", api.AddGroup)
	sysGroup.PUT("group/modify", api.ModifyGroup)
	sysGroup.DELETE("group/delete", api.DeleteGroup)
	sysGroup.POST("group/add/user", api.AddUserToGroup)
	sysGroup.POST("group/delete/user", api.DeleteUserToGroup)
	return nil
}

func NewGroupApi(c framework.Container) *ApiGroup {
	return &ApiGroup{service: group.NewService(c)}
}

package menu

import (
	"devops-http/app/module/sys/service/menu"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiMenu struct {
	service *menu.Service
}

func Register(r *gin.Engine) error {
	api := NewSysApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})

	// 菜单相关接口
	sysGroup.GET("menu/:id", api.GetMenu)
	sysGroup.GET("menu/role/:id", api.GetMenuByRole)
	sysGroup.GET("menu/user", api.GetMenuByUser)
	sysGroup.POST("menu/add/role", api.AddMenuToRole)
	sysGroup.GET("menu/get/tree", api.GetBaseMenuTree)
	sysGroup.POST("menu/list", api.ListMenu)
	sysGroup.POST("menu/add", api.AddMenu)
	sysGroup.PUT("menu/modify", api.UpdateBaseMenu)
	sysGroup.DELETE("menu/delete", api.DeleteBaseMenu)

	return nil
}

func NewSysApi(c framework.Container) *ApiMenu {
	return &ApiMenu{service: menu.NewService(c)}
}

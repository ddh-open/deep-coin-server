package path

import (
	"devops-http/app/module/sys/service/path"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/gin"
)

type ApiPath struct {
	service *path.Service
}

func Register(r *gin.Engine) error {
	api := NewSysApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})

	// api相关接口
	sysGroup.POST("api/get", api.GetApiById)
	sysGroup.POST("api/list", api.GetApiList)
	sysGroup.POST("api/add", api.CreateApi)
	sysGroup.PUT("api/modify", api.UpdateApi)
	sysGroup.DELETE("api/delete", api.DeleteApi)
	sysGroup.DELETE("api/deletes", api.DeleteApisByIds)
	sysGroup.GET("api/all", api.GetAllApis)
	sysGroup.POST("api/role", api.RelativeToRole)
	return nil
}

func NewSysApi(c framework.Container) *ApiPath {
	return &ApiPath{service: path.NewService(c)}
}

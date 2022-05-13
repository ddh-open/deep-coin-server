package tencent

import (
	"devops-http/app/module/third/service/tencent"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiTencent struct {
	service *tencent.Service
}

func Register(r *gin.Engine) error {
	api := NewThirdApi(r.GetContainer())
	userGroup := r.Group("/third/", func(c *gin.Context) {
	})
	// 腾讯云相关接口
	userGroup.POST("tencent/resource/list", api.TencentListResource) // 获取腾讯云的资源（主机，数据库，redis and soon）
	return nil
}

func NewThirdApi(c framework.Container) *ApiTencent {
	return &ApiTencent{service: tencent.NewService(c)}
}

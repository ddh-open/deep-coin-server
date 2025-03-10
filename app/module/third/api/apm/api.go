package apm

import (
	"devops-http/app/module/third/service/apm"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiApm struct {
	service *apm.Service
}

func Register(r *gin.Engine) error {
	api := NewApiApm(r.GetContainer())
	userGroup := r.Group("/third/", func(c *gin.Context) {
	})
	// apm相关接口
	userGroup.POST("tencent/apm/addMerchantApm", api.AddMerchantApm)
	return nil
}

func NewApiApm(c framework.Container) *ApiApm {
	return &ApiApm{service: apm.NewService(c)}
}

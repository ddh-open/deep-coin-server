package cls

import (
	"devops-http/app/module/third/service/cls"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiCls struct {
	service *cls.Service
}

func Register(r *gin.Engine) error {
	api := NewThirdApi(r.GetContainer())
	userGroup := r.Group("/third/", func(c *gin.Context) {
	})
	// cls日志相关接口
	// topic相关接口
	userGroup.POST("tencent/cls/topic/addMerchantClsLogTopic", api.AddMerchantClsLogTopic)
	userGroup.POST("tencent/cls/topic/deleteMerchantLog", api.DeleteMerchantLog)
	return nil
}

func NewThirdApi(c framework.Container) *ApiCls {
	return &ApiCls{service: cls.NewService(c)}
}

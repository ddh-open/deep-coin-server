package domain

import (
	"devops-http/app/module/sys/service/domain"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiDomain struct {
	service *domain.Service
}

func Register(r *gin.Engine) error {
	api := NewDomainApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})
	// 域相关接口
	sysGroup.GET("domain/:id", api.GetDomains)
	sysGroup.POST("domain/list", api.ListDomains)
	sysGroup.POST("domain/add", api.AddDomain)
	sysGroup.PUT("domain/modify", api.ModifyDomain)
	sysGroup.DELETE("domain/delete", api.DeleteDomain)
	return nil
}

func NewDomainApi(c framework.Container) *ApiDomain {
	return &ApiDomain{service: domain.NewService(c)}
}

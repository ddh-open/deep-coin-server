package icon

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/service/icon"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiIcon struct {
	service *icon.Service
}

func Register(r *gin.Engine) error {
	api := NewSysApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})
	sysGroup.POST("icon/list", api.GetList)
	return nil
}

func NewSysApi(c framework.Container) *ApiIcon {
	return &ApiIcon{service: icon.NewService(c)}
}

// GetList godoc
// @Summary 获得图标接口
// @Security ApiKeyAuth
// @Description 获得图标接口
// @accept application/json
// @Produce application/json
// @Param data body request.SearchIconParams true "页数，页大小，筛选条件"
// @Tags Menu
// @Success 200 {object}  response.Response
// @Router /sys/icon/list [post]
func (a *ApiIcon) GetList(c *gin.Context) {
	var param request.SearchIconParams
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := a.service.GetList(param)
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

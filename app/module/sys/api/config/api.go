package config

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/sys"
	configModel "devops-http/app/module/sys/model/config"
	"devops-http/app/module/sys/service/config"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/gin"
	"go.uber.org/zap"
)

type ApiConfig struct {
	service *config.Service
}

func Register(r *gin.Engine) error {
	api := NewConfigApi(r.GetContainer())
	sysGroup := r.Group("/sys/config", func(c *gin.Context) {
	})
	// 域相关接口
	sysGroup.POST("add", api.Create)
	sysGroup.POST("list", api.List)
	sysGroup.PUT("modify", api.Update)
	sysGroup.DELETE("delete", api.Delete)
	return nil
}

func NewConfigApi(c framework.Container) *ApiConfig {
	return &ApiConfig{service: config.NewService(c)}
}

// Create
// @Tags Config
// @Summary 创建config
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body configModel.DevopsSysConfig true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "创建基础api"
// @Router /sys/config/add [post]
func (a *ApiConfig) Create(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req configModel.DevopsSysConfig
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
		response.FailWithMessage("参数解析错误!", c)
		return
	}
	if err = a.service.Create(req); err != nil {
		logGet.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// Delete
// @Tags Config
// @Summary 删除config
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys.RequestById true "ID"
// @Success 200 {object} response.Response{msg=string} "删除api"
// @Router /sys/config/delete [delete]
func (a *ApiConfig) Delete(c *gin.Context) {
	logGet := c.MustMakeLog()
	var ids sys.RequestById
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
		response.FailWithMessage("参数解析错误!", c)
		return
	}
	if err := a.service.Delete(ids); err != nil {
		logGet.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// Update
// @Tags Config
// @Summary 修改config
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body configModel.DevopsSysConfig true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "修改基础api"
// @Router /sys/config/modify [put]
func (a *ApiConfig) Update(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req configModel.DevopsSysConfig
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
		response.FailWithMessage("参数解析错误!", c)
		return
	}
	if err := a.service.Update(req); err != nil {
		logGet.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// List
// @Tags Config
// @Summary 分页获取config列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SearchConfigParams true "分页获取API列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router /sys/config/list [post]
func (a *ApiConfig) List(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req request.SearchConfigParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
	}
	if err, list, total := a.service.List(req); err != nil {
		logGet.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, "获取成功", c)
	}
}

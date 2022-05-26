package application

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	applicationModel "devops-http/app/module/cluster/model/application"
	"devops-http/app/module/cluster/service/application"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiApplication struct {
	service *application.Service
}

func NewApplicationApi(c framework.Container) *ApiApplication {
	return &ApiApplication{service: application.NewService(c)}
}

// List godoc
// @Summary 获得应用列表接口
// @Security ApiKeyAuth
// @Description 获得应用列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.SearchApplicationParams true "页数，页大小，筛选条件"
// @Tags Menu
// @Success 200 {object}  response.Response
// @Router /cluster/application/list [get]
func (a *ApiApplication) List(c *gin.Context) {
	var param request.SearchApplicationParams
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := a.service.List(param)
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = &result
	c.DJson(res)
}

// GetApplication godoc
// @Summary 获得单个应用接口
// @Security ApiKeyAuth
// @Description 获得单个应用接口
// @accept application/json
// @Produce application/json
// @Param id path int true "应用id"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/{id} [get]
func (a *ApiApplication) GetApplication(c *gin.Context) {
	applicationId := c.Param("id")
	result, err := a.service.GetApplicationById(applicationId)
	res := response.Response{Code: 1, Msg: "查询成功", Data: result}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// AddApplication godoc
// @Summary 新增应用接口
// @Security ApiKeyAuth
// @Description 新增应用接口
// @accept application/json
// @Produce application/json
// @Param data body applicationModel.DevopsClusterApplication true "应用"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/add [post]
func (a *ApiApplication) AddApplication(c *gin.Context) {
	var req applicationModel.DevopsClusterApplication
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}

	err = a.service.Save(&req)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// Modify godoc
// @Summary 修改应用接口
// @Security ApiKeyAuth
// @Description 修改应用接口
// @accept application/json
// @Produce application/json
// @Param data body role.DevopsSysRoleEntity true "角色"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/modify [put]
func (a *ApiApplication) Modify(c *gin.Context) {
	var req applicationModel.DevopsClusterApplication
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "修改成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.Modify(&req)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// AddApplicationConfig godoc
// @Summary 新增应用接口
// @Security ApiKeyAuth
// @Description 新增应用接口
// @accept application/json
// @Produce application/json
// @Param data body applicationModel.DevopsClusterApplication true "应用"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/add/config [post]
func (a *ApiApplication) AddApplicationConfig(c *gin.Context) {
	var req applicationModel.DevopsClusterApplication
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}

	err = a.service.AddConfig(&req)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// ModifyConfig godoc
// @Summary 修改应用配置接口
// @Security ApiKeyAuth
// @Description 修改应用配置接口
// @accept application/json
// @Produce application/json
// @Param data body applicationModel.DevopsClusterApplication true "修改应用配置"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/modify/config [put]
func (a *ApiApplication) ModifyConfig(c *gin.Context) {
	var req applicationModel.DevopsClusterApplication
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "修改成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.ModifyConfig(&req)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// DeleteConfig godoc
// @Summary 删除应用配置接口
// @Security ApiKeyAuth
// @Description 删除应用配置接口
// @accept application/json
// @Produce application/json
// @Param data body applicationModel.DevopsClusterApplication true "删除应用配置"
// @Tags application
// @Success 200 {object}  response.Response
// @Router /cluster/application/delete/config [delete]
func (a *ApiApplication) DeleteConfig(c *gin.Context) {
	var req applicationModel.DevopsClusterApplication
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.DeleteConfig(&req)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// Delete godoc
// @Summary 删除应用接口
// @Security ApiKeyAuth
// @Description 删除应用接口
// @accept application/json
// @Produce application/json
// @Param id path int true "应用id"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /cluster/application/delete/{id} [delete]
func (a *ApiApplication) Delete(c *gin.Context) {
	applicationId := c.Param("id")
	res := response.Response{Code: 1, Msg: "删除成功"}
	err := a.service.Delete(applicationId)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

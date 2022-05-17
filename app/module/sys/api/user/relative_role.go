package user

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/utils"
	"devops-http/framework/gin"
)

// UserRelativeRole godoc
// @Summary 获得用户角色接口
// @Security ApiKeyAuth
// @Description 获得用户角色接口
// @accept application/json
// @Produce application/json
// @Param id path int true "用户id"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/relative/roles/{id} [get]
func (api *ApiUser) UserRelativeRole(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	id := c.Param("id")
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	result, err := api.service.GetRolesByUserId(id, userToken.CurrentDomain, cabin)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = result
	c.DJson(res)
}

// UserGetApis godoc
// @Summary 获得用户api权限接口
// @Security ApiKeyAuth
// @Description 获得用户api权限接口
// @accept application/json
// @Produce application/json
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/relative/apis [get]
func (api *ApiUser) UserGetApis(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "获取成功", Data: nil}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	result, err := api.service.GetUserApis(userToken, cabin)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = map[string]interface{}{"list": result}
	c.DJson(res)
}

// UserRelativeRoleAdd godoc
// @Summary 用户关联角色接口
// @Security ApiKeyAuth
// @Description 用户关联角色接口
// @accept application/json
// @Produce application/json
// @Param data body request.UserRelativeRoleRequest true "关联参数"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/relative/roles/add [post]
func (api *ApiUser) UserRelativeRoleAdd(c *gin.Context) {
	logger := c.MustMakeLog()
	param := request.UserRelativeRoleRequest{}
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	err = api.service.RelativeRolesToUser(param, userToken.CurrentDomain, cabin)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// UserRelativeRoleDelete godoc
// @Summary 删除用户角色接口
// @Security ApiKeyAuth
// @Description 删除用户角色接口
// @accept application/json
// @Produce application/json
// @Param data body request.UserRelativeRoleRequest true "关联参数"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/relative/roles/delete [delete]
func (api *ApiUser) UserRelativeRoleDelete(c *gin.Context) {
	logger := c.MustMakeLog()
	param := request.UserRelativeRoleRequest{}
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	err = api.service.DeleteRelativeRolesToUser(param, userToken.CurrentDomain, cabin)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

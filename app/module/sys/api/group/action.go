package group

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
)

// GetGroups godoc
// @Summary 获得分组接口
// @Security ApiKeyAuth
// @Description 获得分组接口
// @accept application/json
// @Produce application/json
// @Param id path int true "分组id"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/{id} [get]
func (a *ApiGroup) GetGroups(c *gin.Context) {
	id := c.Param("id")
	result, err := a.service.GetGroupById(id)
	res := response.Response{Code: 1, Msg: "查询成功", Data: result}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// ListGroups godoc
// @Summary 获得分组列表接口
// @Security ApiKeyAuth
// @Description 获得分组列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.SearchGroupParams true "页数，页大小，筛选条件"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/list [post]
func (a *ApiGroup) ListGroups(c *gin.Context) {
	var param request.SearchGroupParams
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := a.service.GetGroupList(param)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

// TreeGroups godoc
// @Summary 获得分组树形结构接口
// @Security ApiKeyAuth
// @Description 获得分组树形结构接口
// @accept application/json
// @Produce application/json
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/tree [get]
func (a *ApiGroup) TreeGroups(c *gin.Context) {
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	result, err := a.service.GetGroupTree()
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

// AddGroup godoc
// @Summary 新增分组接口
// @Security ApiKeyAuth
// @Description 新增分组接口
// @accept application/json
// @Produce application/json
// @Param data body base.DevopsSysGroup true "分组"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/add [post]
func (a *ApiGroup) AddGroup(c *gin.Context) {
	var req base.DevopsSysGroup
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.AddGroup(req)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// ModifyGroup godoc
// @Summary 修改分组接口
// @Security ApiKeyAuth
// @Description 修改分组接口
// @accept application/json
// @Produce application/json
// @Param data body base.DevopsSysGroup true "分组"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/modify [post]
func (a *ApiGroup) ModifyGroup(c *gin.Context) {
	var req base.DevopsSysGroup
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "修改成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.ModifyGroup(req)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// DeleteGroup godoc
// @Summary 删除分组接口
// @Security ApiKeyAuth
// @Description 删除分组接口
// @accept application/json
// @Produce application/json
// @Param data body request.ReqById true "分组"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/delete [delete]
func (a *ApiGroup) DeleteGroup(c *gin.Context) {
	var req request.ReqById
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.DeleteGroup(req)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// AddUserToGroup godoc
// @Summary 给分组新增用户
// @Security ApiKeyAuth
// @Description 给分组新增用户
// @accept application/json
// @Produce application/json
// @Param data body request.GroupRelativeUserRequest true "给分组新增用户"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/group/add/user [post]
func (a *ApiGroup) AddUserToGroup(c *gin.Context) {
	var req request.GroupRelativeUserRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.AddUserToGroup(req)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// DeleteUserToGroup godoc
// @Summary 给分组删除用户
// @Security ApiKeyAuth
// @Description 给分组删除用户
// @accept application/json
// @Produce application/json
// @Param data body request.GroupRelativeUserRequest true "给分组新增用户"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/group/delete/user [post]
func (a *ApiGroup) DeleteUserToGroup(c *gin.Context) {
	var req request.GroupRelativeUserRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.DeleteUserToGroup(req)
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

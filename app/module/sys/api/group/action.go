package group

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/group"
	"devops-http/framework/gin"
)

//// GetGroupsResource godoc
//// @Summary 获得分组资源接口
//// @Security ApiKeyAuth
//// @Description 获得分组资源接口
//// @accept application/json
//// @Produce application/json
//// @Param name path string true "分组name"
//// @Param domain query string false "域"
//// @Tags Group
//// @Success 200 {object}  response.Response
//// @Router /sys/group/resource/{name} [get]
//func (a *ApiGroup) GetGroupsResource(c *gin.Context) {
//	name := c.Param("name")
//	domain := c.Query("domain")
//	// 从cookie中获取domain
//	if d, err := c.Cookie("domain"); err == nil {
//		domain = d
//	}
//	result, err := a.service.GetGroupsResource(name, domain, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
//	res := response.Response{Code: 1, Msg: "查询成功", Data: result}
//	if err != nil {
//		res.Code = -1
//		res.Msg = err.Error()
//	}
//	c.DJson(res)
//}

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

// AddGroup godoc
// @Summary 新增分组接口
// @Security ApiKeyAuth
// @Description 新增分组接口
// @accept application/json
// @Produce application/json
// @Param data body group.DevopsSysGroup true "分组"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/add [post]
func (a *ApiGroup) AddGroup(c *gin.Context) {
	var req group.DevopsSysGroup
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
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
// @Param data body group.DevopsSysGroup true "分组"
// @Tags Group
// @Success 200 {object}  response.Response
// @Router /sys/group/modify [post]
func (a *ApiGroup) ModifyGroup(c *gin.Context) {
	var req group.DevopsSysGroup
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

// AddResourcesToGroup godoc
// @Summary 给分组新增资源接口
// @Security ApiKeyAuth
// @Description 给分组新增资源接口
// @accept application/json
// @Produce application/json
// @Param data body []request.CabinInReceive true "Ptype为p2 , source 是分组的id，resource 是资源， method 为write或者read,或者owner"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/group/add/resources [post]
func (a *ApiGroup) AddResourcesToGroup(c *gin.Context) {
	param := make([]request.CabinInReceive, 0)
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	err = a.service.AddResourcesToGroup(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
	}
	c.DJson(res)
}

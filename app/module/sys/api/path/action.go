package path

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/sys/model/path"
	"devops-http/framework/gin"
	"go.uber.org/zap"
)

// CreateApi
// @Tags Apis
// @Summary 创建基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body path.DevopsSysApi true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "创建基础api"
// @Router /sys/api/add [post]
func (a *ApiPath) CreateApi(c *gin.Context) {
	logGet := c.MustMakeLog()
	var api path.DevopsSysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
		response.FailWithMessage("参数解析错误!", c)
		return
	}
	if err = a.service.CreateApi(api); err != nil {
		logGet.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// GetApiList
// @Tags Apis
// @Summary 分页获取API列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SearchApiParams true "分页获取API列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router /sys/api/list [post]
func (a *ApiPath) GetApiList(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req request.SearchApiParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
	}
	if err, list, total := a.service.GetAPIInfoList(req, logGet); err != nil {
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

// GetApiById
// @Tags Apis
// @Summary 根据id获取api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ReqById true "根据id获取api"
// @Success 200 {object} response.Response{data=path.DevopsSysApi} "根据id获取api,返回包括api详情"
// @Router /sys/api/get [post]
func (a *ApiPath) GetApiById(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req request.ReqById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	id := c.Param("id")
	err, api := a.service.GetApiById(id)
	if err != nil {
		logGet.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(api, c)
	}
}

// RelativeToRole
// @Tags Apis
// @Summary 给角色关联api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RelativeRoleApisRequest true "给角色关联api 权限"
// @Success 200 {object} response.Response{} "给角色关联api权限"
// @Router /sys/api/role [post]
func (a *ApiPath) RelativeToRole(c *gin.Context) {
	logGet := c.MustMakeLog()
	var req request.RelativeRoleApisRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logGet.Error("参数获取失败!", zap.Error(err))
		response.FailWithMessage("参数获取失败", c)
		return
	}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logGet.Error("用户token 解析失败!", zap.Error(err))
		response.FailWithMessage("用户token 解析失败", c)
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	err = a.service.RelativeApiToRole(req, userToken, cabin)
	if err != nil {
		logGet.Error("关联!", zap.Error(err))
		response.FailWithMessage("关联", c)
	} else {
		response.OkWithMessage("关联成功", c)
	}
}

// UpdateApi
// @Tags Apis
// @Summary 修改基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body path.DevopsSysApi true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Response{msg=string} "修改基础api"
// @Router /sys/api/modify [put]
func (a *ApiPath) UpdateApi(c *gin.Context) {
	logGet := c.MustMakeLog()
	var api path.DevopsSysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		logGet.Error("参数解析错误!", zap.Error(err))
		response.FailWithMessage("参数解析错误!", c)
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	if err := a.service.UpdateApi(api, cabin); err != nil {
		logGet.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// GetAllApis
// @Tags Apis
// @Summary 获取所有的Api 不分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string][]path.DevopsSysApi,msg=string} "获取所有的Api 不分页,返回包括api列表"
// @Router /sys/api/all [get]
func (a *ApiPath) GetAllApis(c *gin.Context) {
	logGet := c.MustMakeLog()
	if err, apis := a.service.GetAllApis(); err != nil {
		logGet.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(apis, "获取成功", c)
	}
}

// DeleteApisByIds
// @Tags Apis
// @Summary 删除选中Api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ReqById true "ID"
// @Success 200 {object} response.Response{msg=string} "删除选中Api"
// @Router /sys/api/delete [delete]
func (a *ApiPath) DeleteApisByIds(c *gin.Context) {
	logGet := c.MustMakeLog()
	var ids request.ReqById
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		logGet.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	if err := a.service.DeleteApisByIds(ids); err != nil {
		logGet.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

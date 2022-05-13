package role

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
)

// GetRoles godoc
// @Summary 获得角色接口
// @Security ApiKeyAuth
// @Description 获得角色接口
// @accept application/json
// @Produce application/json
// @Param id path int true "角色id"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/{id} [get]
func (a *ApiRole) GetRoles(c *gin.Context) {
	roleId := c.Param("id")
	result, err := a.service.GetRoleById(roleId, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	res := response.Response{Code: 1, Msg: "查询成功", Data: result}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// ListRoles godoc
// @Summary 获得角色列表接口
// @Security ApiKeyAuth
// @Description 获得角色列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.PageRequest true "页数，页大小，筛选条件"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/list [post]
func (a *ApiRole) ListRoles(c *gin.Context) {
	var param request.PageRequest
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	result, err := a.service.GetRoleList(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		return
	}
	res.Data = result
	c.DJson(res)
}

// AddRole godoc
// @Summary 新增角色接口
// @Security ApiKeyAuth
// @Description 新增角色接口
// @accept application/json
// @Produce application/json
// @Param data body role.DevopsSysRole true "角色"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/add [post]
func (a *ApiRole) AddRole(c *gin.Context) {
	param := make(map[string]interface{}, 0)
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.AddRole(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// ModifyRole godoc
// @Summary 修改角色接口
// @Security ApiKeyAuth
// @Description 修改角色接口
// @accept application/json
// @Produce application/json
// @Param data body role.DevopsSysRole true "角色"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/modify [post]
func (a *ApiRole) ModifyRole(c *gin.Context) {
	var param map[string]interface{}
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "修改成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.ModifyRole(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// DeleteRole godoc
// @Summary 删除角色接口
// @Security ApiKeyAuth
// @Description 删除角色接口
// @accept application/json
// @Produce application/json
// @Param ids body request.DataDelete true "角色ids"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/delete [delete]
func (a *ApiRole) DeleteRole(c *gin.Context) {
	var req request.DataDelete
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.DeleteRole(req.Ids, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// AddResourcesToRole godoc
// @Summary 给角色新增权限接口
// @Security ApiKeyAuth
// @Description 给角色新增权限接口
// @accept application/json
// @Produce application/json
// @Param data body []request.CabinInReceive true "新增api权限Ptype为p； 新增菜单权限Ptype为p3 , source 是角色的id，resource 是资源"
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/add/resources [post]
func (a *ApiRole) AddResourcesToRole(c *gin.Context) {
	param := make([]request.CabinInReceive, 0)
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "新增成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.AddResourcesToRole(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

// CopyRole godoc
// @Summary 复制角色接口
// @Security ApiKeyAuth
// @Description 复制角色接口
// @accept application/json
// @Produce application/json
// @Param data body request.CopyRoleParams true "copyId "
// @Tags Role
// @Success 200 {object}  response.Response
// @Router /sys/roles/copy [post]
func (a *ApiRole) CopyRole(c *gin.Context) {
	var req request.CopyRoleParams
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "复制成功"}
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		c.DJson(res)
		return
	}
	err = a.service.CopyRole(req, c.MustMake(contract.KeyCaBin).(contract.Cabin))
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
	}
	c.DJson(res)
}

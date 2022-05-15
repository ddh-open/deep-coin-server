package menu

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/sys"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/sys/model/menu"
	"devops-http/framework/gin"
	"github.com/pkg/errors"
)

// AddMenu godoc
// @Summary 新增菜单接口
// @Security ApiKeyAuth
// @Description 新增菜单接口
// @accept application/json
// @Produce application/json
// @Param data body menu.DevopsSysMenuEntity true "菜单"
// @Tags Menu
// @Success 200 {object}  response.Response
// @Router /sys/menu/add [post]
func (a *ApiMenu) AddMenu(c *gin.Context) {
	logger := c.MustMakeLog()
	var menuData menu.DevopsSysMenuEntity
	res := response.Response{Code: 1, Msg: "新增成功"}
	err := c.ShouldBindJSON(&menuData)
	if err != nil {
		res.Msg = "添加失败: " + err.Error()
		logger.Error(res.Msg)
		res.Code = -1
		c.DJson(res)
		return
	}
	if err = a.service.AddBaseMenu(menuData); err != nil {
		res.Msg = "添加失败: " + err.Error()
		logger.Error(res.Msg)
		res.Code = -1
	}
	c.DJson(res)
}

// DeleteBaseMenu @Tags Menu
// @Summary 删除菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys.DeleteById true "菜单id"
// @Tags Menu
// @Success 200 {object} response.Response{msg=string} "删除菜单"
// @Router /sys/menu/delete [delete]
func (a *ApiMenu) DeleteBaseMenu(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "删除成功"}
	var param sys.DeleteById
	err := c.ShouldBindJSON(&param)
	if err != nil {
		res.Msg = "删除失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
		c.DJson(res)
	}
	if err := a.service.DeleteBaseMenu(param); err != nil {
		res.Msg = "删除失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
	}
	c.DJson(res)
}

// UpdateBaseMenu @Tags Menu
// @Summary 更新菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body menu.DevopsSysMenu true "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Tags Menu
// @Success 200 {object} response.Response{msg=string} "更新菜单"
// @Router /sys/menu/modify [put]
func (a *ApiMenu) UpdateBaseMenu(c *gin.Context) {
	var menuData menu.DevopsSysMenuEntity
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "更新成功"}
	_ = c.ShouldBindJSON(&menuData)
	if err := a.service.UpdateBaseMenu(menuData); err != nil {
		res.Msg = "更新失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
	}
	c.DJson(res)
}

// GetMenu godoc
// @Summary 获得菜单接口
// @Security ApiKeyAuth
// @Description 获得菜单接口
// @accept application/json
// @Produce application/json
// @Param id path int true "菜单id"
// @Tags Menu
// @Success 200 {object}  response.Response
// @Router /sys/menu/{id} [get]
func (a *ApiMenu) GetMenu(c *gin.Context) {
	id := c.Param("id")
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "获取成功"}
	if err, menuData := a.service.GetBaseMenuById(id); err != nil {
		res.Msg = "获取失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
	} else {
		res.Data = menuData
	}
	c.DJson(res)
}

// ListMenu godoc
// @Summary 获得菜单列表接口
// @Security ApiKeyAuth
// @Description 获得菜单列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.SearchMenusParams true "页数，页大小，筛选条件"
// @Tags Menu
// @Success 200 {object}  response.Response
// @Router /sys/menu/list [post]
func (a *ApiMenu) ListMenu(c *gin.Context) {
	var param request.SearchMenusParams
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := a.service.GetList(param)
	if err != nil {
		res.Code = -1
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

// GetBaseMenuTree @Tags GetBaseMenuTree
// @Summary 获取用户动态路由
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} []menu.DevopsSysMenu
// @Tags Menu
// @Router /sys/menu/get/tree [get]
func (a *ApiMenu) GetBaseMenuTree(c *gin.Context) {
	res := response.Response{Code: 1, Msg: "获取成功"}
	logger := c.MustMakeLog()
	if err, menus := a.service.GetBaseMenuTree(); err != nil {
		res.Msg = "获取失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
		c.DJson(res)
		return
	} else {
		res.Data = map[string]interface{}{"list": menus}
		c.DJson(res)
	}
}

// AddMenuToRole @Tags AddMenuToRole
// @Summary 增加menu和角色关联关系
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RelativeRoleMenuRequest true "角色ID"
// @Tags Menu
// @Success 200 {object} response.Response{msg=string} "增加menu和角色关联关系"
// @Router /sys/menu/add/role [post]
func (a *ApiMenu) AddMenuToRole(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{Code: -1, Msg: "添加成功"}
	var req request.RelativeRoleMenuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.Msg = errors.Errorf("参数解析错误：%s", err).Error()
		logger.Error(res.Msg)
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
	if err := a.service.AddMenuToRole(userToken, req, cabin); err != nil {
		res.Msg = "添加失败: " + err.Error()
		logger.Error(res.Msg)
	}
	res.Code = 1
	c.DJson(res)
}

// GetMenuByRole @Tags GetMenuByRole
// @Summary 获取指定角色menu
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "角色ID"
// @Tags Menu
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取指定角色menu"
// @Router /sys/menu/role/{id} [get]
func (a *ApiMenu) GetMenuByRole(c *gin.Context) {
	roleId := c.Param("id")
	logger := c.MustMakeLog()
	res := response.Response{Code: -1, Msg: "获取失败"}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	if err, menus := a.service.GetMenuByRole(userToken, roleId, cabin); err != nil {
		res.Msg = "获取失败: " + err.Error()
		logger.Error(res.Msg)
	} else {
		res.Data = menus
	}
	res.Code = 1
	res.Msg = "获取成功"
	c.DJson(res)
}

// GetMenuByUser @Tags GetMenuByUser
// @Summary 获取用户的menu
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Tags Menu
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取指定角色menu"
// @Router /sys/menu/user [get]
func (a *ApiMenu) GetMenuByUser(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{Code: 1, Msg: "获取成功", Data: nil}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	if err, menus := a.service.GetMenuByUser(userToken, cabin); err != nil {
		res.Msg = "获取失败: " + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
	} else {
		res.Data = map[string]interface{}{"list": menus}
	}
	c.DJson(res)
}

package user

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/sys"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/sys/model/user"
	"devops-http/framework/gin"
)

// Login godoc
// @Summary 用户登录接口
// @Security ApiKeyAuth
// @Description 用户登录接口
// @accept application/json
// @Produce application/json
// @Param data body sys.LoginRequest true "用户名，密码，账户类型"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/login [post]
func (api *ApiUser) Login(c *gin.Context) {
	logger := c.MustMakeLog()
	var req sys.LoginRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	jwt := c.MustMake(contract.JWT).(contract.JWTService)
	data, err := api.service.Login(req, jwt)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = data
	c.DJson(res)
}

// Modify godoc
// @Summary 用户修改接口
// @Security ApiKeyAuth
// @Description 用户修改接口
// @accept application/json
// @Produce application/json
// @Param data body user.DevopsSysUserEntity true "用户的相关字段， 用户的旧密码"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/modify [put]
func (api *ApiUser) Modify(c *gin.Context) {
	logger := c.MustMakeLog()
	var request user.DevopsSysUserEntity
	err := c.ShouldBindJSON(&request)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	data, err := api.service.Modify(request, cabin)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = data
	c.DJson(res)
}

// Add godoc
// @Summary 用户新增接口
// @Security ApiKeyAuth
// @Description 用户新增接口
// @accept application/json
// @Produce application/json
// @Param data body user.DevopsSysUserEntity true "用户的相关字段， 用户的旧密码"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/add [post]
func (api *ApiUser) Add(c *gin.Context) {
	logger := c.MustMakeLog()
	var request user.DevopsSysUserEntity
	err := c.ShouldBindJSON(&request)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	ldap := c.MustMake(contract.KeyLdap).(contract.Ldap)
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	data, err := api.service.Add(request, ldap, cabin)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = data
	c.DJson(res)
}

// Delete godoc
// @Summary 用户删除接口
// @Security ApiKeyAuth
// @Description 用户删除接口
// @accept application/json
// @Produce application/json
// @Param ids body request.DataDelete true "传用户id, 多个以, 隔开"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/delete [delete]
func (api *ApiUser) Delete(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.DataDelete
	err := c.ShouldBindJSON(&req)
	res := response.Response{Code: 1, Msg: "删除成功"}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	cabin := c.MustMake(contract.KeyCaBin).(contract.Cabin)
	err = api.service.Delete(req.Ids, cabin)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// ChangePassword godoc
// @Summary 用户修改密码接口
// @Security ApiKeyAuth
// @Description 用户修改密码接口
// @accept application/json
// @Produce application/json
// @Param data body sys.ChangePasswordRequest true "用户名，原密码， 新密码，账户类型"
// @Tags User
// @Success 200 {object}  response.Response
// @Router /user/changePassword [post]
func (api *ApiUser) ChangePassword(c *gin.Context) {
	logger := c.MustMakeLog()
	var req sys.ChangePasswordRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	ldap := c.MustMake(contract.KeyLdap).(contract.Ldap)
	err = api.service.ChangePassword(req, ldap)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	c.DJson(res)
}

// Register godoc
// @Summary 用户注册接口
// @Security ApiKeyAuth
// @Description 用户注册接口
// @Produce  json
// @Tags User
// @Success 200 {object} response.Response
// @Router /user/register [post]
func (api *ApiUser) Register(c *gin.Context) {
	logger := c.MustMakeLog()
	var request sys.LoginRequest
	err := c.ShouldBindJSON(&request)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	jwt := c.MustMake(contract.JWT).(contract.JWTService)
	data, err := api.service.Login(request, jwt)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = data
	c.DJson(res)
}

// Logout godoc
// @Summary 用户退出接口
// @Security ApiKeyAuth
// @Description 用户退出接口
// @Produce  json
// @Tags User
// @Success 200 {object} response.Response
// @Router /user/logout [post]
func (api *ApiUser) Logout(c *gin.Context) {
	logger := c.MustMakeLog()
	var req sys.LoginRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	jwt := c.MustMake(contract.JWT).(contract.JWTService)
	data, err := api.service.Login(req, jwt)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	}
	res.Data = data
	c.DJson(res)
}

// GetUserInfo godoc
// @Summary 获取用户详情的接口
// @Security ApiKeyAuth
// @Description 根据token获取id
// @Produce  json
// @Tags User
// @Success 200 {object} user.DevopsSysUserView
// @Router /user/info [get]
func (api *ApiUser) GetUserInfo(c *gin.Context) {
	logger := c.MustMakeLog()
	res := response.Response{
		Code: -1,
		Data: nil,
	}
	userToken, err := utils.ParseToken(c)
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
		return
	}
	data, err := api.service.UserInfo(userToken, c.MustMake(contract.KeyCaBin).(contract.Cabin), append([]interface{}{}, "id = ?", userToken.Id))
	if err != nil {
		logger.Error(err.Error())
		res.Msg = err.Error()
	}
	res.Code = 1
	res.Data = data
	c.DJson(res)
}

// UserList godoc
// @Summary 获取用户列表的接口
// @Security ApiKeyAuth
// @Description 根据参数获取用户列表
// @Produce  json
// @Param data body request.PageRequest true "分页查询"
// @Tags User
// @Success 200 {object} response.PageResult
// @Router /user/list [post]
func (api *ApiUser) UserList(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.PageRequest
	err := c.ShouldBindJSON(&req)
	res := response.Response{
		Code: 1,
		Data: nil,
	}
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
		return
	}
	data, err := api.service.UserList(c.MustMake(contract.KeyCaBin).(contract.Cabin), req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
		return
	}
	res.Data = data
	c.DJson(res)
}

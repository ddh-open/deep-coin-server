package process

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/workflow/service/process"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type instProcess struct {
	service *process.Service
}

func NewInstProcess(c framework.Container) *instProcess {
	return &instProcess{service: process.NewService(c)}
}

// Start godoc
// @Summary 开启工作流实例
// @Security ApiKeyAuth
// @Description 开启工作流实例
// @accept application/json
// @Produce application/json
// @Param data body request.ReceiverProcess true "模型ReceiverProcess 中的数据"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/process/start [post]
func (api *instProcess) Start(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.ReceiverProcess
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
	id, err := api.service.StartProcess(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = id
	}
	c.DJson(res)
}

// FindProcInst godoc
// @Summary 查询到我审批的流程实例
// @Security ApiKeyAuth
// @Description 查询到我审批的流程实例
// @accept application/json
// @Produce application/json
// @Param data body request.ProcessPageReceiver true "查询轮到自己审批的流程"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/process/inst [post]
func (api *instProcess) FindProcInst(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.ProcessPageReceiver
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
	data, err := api.service.FindProcList(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

// FindProcInstMyself godoc
// @Summary 查询到我发起的流程实例
// @Security ApiKeyAuth
// @Description 查询到我发起的流程实例
// @accept application/json
// @Produce application/json
// @Param data body request.ProcessPageReceiver true "查询到我发起的流程实例"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/process/inst/myself [post]
func (api *instProcess) FindProcInstMyself(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.ProcessPageReceiver
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
	data, err := api.service.FindProcList(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

// FindProcNotifyInst godoc
// @Summary 查询抄送我的流程实例
// @Security ApiKeyAuth
// @Description 查询到我发起的流程实例
// @accept application/json
// @Produce application/json
// @Param data body request.ProcessPageReceiver true "查询到我发起的流程实例"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/process/inst/notify [post]
func (api *instProcess) FindProcNotifyInst(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.ProcessPageReceiver
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
	data, err := api.service.FindProcNotify(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

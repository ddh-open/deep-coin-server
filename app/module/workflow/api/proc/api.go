package proc

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/workflow"
	"devops-http/app/module/workflow/service/node"
	"devops-http/app/module/workflow/service/proc"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/gin"
)

type defProc struct {
	service *proc.Service
}

func NewDefProc(c framework.Container) *defProc {
	node.NewService(c)
	return &defProc{service: proc.NewService(c)}
}

// Save godoc
// @Summary 定义工作流接口
// @Security ApiKeyAuth
// @Description 定义工作流接口
// @accept application/json
// @Produce application/json
// @Param data body request.ProcRequest true "流程名， 用户id， 用户所在的公司，定义流程json串"
// @Tags proc
// @Success 200 {object}  response.Response
// @Router /v1/workflow/proc/save [post]
func (api *defProc) Save(c *gin.Context) {
	//c.GetUser()
	logger := c.MustMakeLog()
	var req request.ProcRequest
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
	id, err := api.service.SaveProc(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = id
	}
	c.DJson(res)
}

// Delete godoc
// @Summary 删除工作流接口
// @Security ApiKeyAuth
// @Description 删除工作流接口
// @accept application/json
// @Produce application/json
// @Param id path int true "工作流id"
// @Tags proc
// @Success 200 {object}  response.Response
// @Router /v1/workflow/proc/delete/{id} [delete]
func (api *defProc) Delete(c *gin.Context) {
	logger := c.MustMakeLog()
	id := c.Param("id")
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	err := api.service.DeleteProc(id)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Msg = "删除失败"
	}
	c.DJson(res)
}

// List godoc
// @Summary 查询定义的工作流列表
// @Security ApiKeyAuth
// @Description 查询定义的工作流列表
// @accept application/json
// @Produce application/json
// @Param data body request.ProcPageReceiver true "page  pageSize  filter 筛选条件"
// @Tags proc
// @Success 200 {object}  response.Response
// @Router /v1/workflow/proc/list [post]
func (api *defProc) List(c *gin.Context) {
	//c.GetUser()
	logger := c.MustMakeLog()
	var req workflow.ProcPageReceiver
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
	list, err := api.service.ListProc(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = list
	}
	c.DJson(res)
}

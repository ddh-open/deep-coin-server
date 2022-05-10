package task

import (
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/workflow"
	"devops-http/app/module/workflow/service/task"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/gin"
)

type taskApi struct {
	service *task.Service
}

func NewTaskApi(c framework.Container) *taskApi {
	return &taskApi{service: task.NewService(c)}
}

// CompleteTask godoc
// @Summary 审批
// @Security ApiKeyAuth
// @Description 审批
// @accept application/json
// @Produce application/json
// @Param data body workflow.TaskReceiver true "完成审批所需的参数"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/task/complete [post]
func (api *taskApi) CompleteTask(c *gin.Context) {
	logger := c.MustMakeLog()
	var req workflow.TaskReceiver
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
	data, err := api.service.CompleteTask(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

// WithDrawTask godoc
// @Summary 撤回
// @Security ApiKeyAuth
// @Description 撤回
// @accept application/json
// @Produce application/json
// @Param data body workflow.TaskReceiver true "撤回所需参数"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/task/complete [post]
func (api *taskApi) WithDrawTask(c *gin.Context) {
	logger := c.MustMakeLog()
	var req workflow.TaskReceiver
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
	data, err := api.service.WithDrawTask(&req)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

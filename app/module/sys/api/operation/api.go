package operation

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/service/operation"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiOperation struct {
	service *operation.Service
}

func Register(r *gin.Engine) error {
	api := NewOperationApi(r.GetContainer())
	sysGroup := r.Group("/sys/", func(c *gin.Context) {
	})
	// 操作记录查询
	sysGroup.GET("operation/:id", api.GetDetail)
	sysGroup.POST("operation/list", api.GetList)
	return nil
}

func NewOperationApi(c framework.Container) *ApiOperation {
	return &ApiOperation{service: operation.NewService(c)}
}

// GetDetail godoc
// @Summary 获得操作记录接口
// @Security ApiKeyAuth
// @Description 获得操作记录接口
// @accept application/json
// @Produce application/json
// @Param id path int true "操作记录的id"
// @Tags Operation
// @Success 200 {object}  response.Response
// @Router /sys/operation/{id} [get]
func (a *ApiOperation) GetDetail(c *gin.Context) {
	logger := c.MustMakeLog()
	id := c.Param("id")
	res := response.Response{
		Code: 1,
		Data: nil,
		Msg:  "查询成功",
	}
	data, err := a.service.GetDetailById(id)
	if err != nil {
		res.Msg = "查询出错：" + err.Error()
		logger.Error(res.Msg)
		res.Code = -1
		c.DJson(res)
		return
	}
	res.Data = data
	c.DJson(res)
}

// GetList godoc
// @Summary 获得操作记录列表接口
// @Security ApiKeyAuth
// @Description 获得操作记录列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.SearchLogsParams true "页数，页大小，筛选条件"
// @Tags Operation
// @Success 200 {object}  []operation.DevopsSysOperationRecord
// @Router /sys/operation/list [post]
func (a *ApiOperation) GetList(c *gin.Context) {
	logger := c.MustMakeLog()
	var req request.SearchLogsParams
	err := c.ShouldBindJSON(&req)
	res := response.Response{
		Code: 1,
		Data: nil,
		Msg:  "查询成功",
	}
	if err != nil {
		res.Msg = "参数解析出错：" + err.Error()
		res.Code = -1
		logger.Error(res.Msg)
		c.DJson(res)
		return
	}
	data, err := a.service.GetList(req)
	if err != nil {
		res.Msg = "查询出错：" + err.Error()
		logger.Error(res.Msg)
		res.Code = -1
		c.DJson(res)
		return
	}
	res.Data = data
	c.DJson(res)
}

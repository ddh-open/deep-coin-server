package cls

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
)

// AddMerchantClsLogTopic godoc
// @Summary 新增商户的日志主题和日志集
// @Security ApiKeyAuth
// @Description 新增商户的日志主题和日志集
// @accept application/json
// @Produce application/json
// @Param data body request.AddMerchantApmRequest true "商户名称，商户id，名称空间"
// @Tags ThirdTencent
// @Success 200 {object}  response.Response
// @Router /third/tencent/cls/topic/addMerchantClsLogTopic [post]
func (t *ApiCls) AddMerchantClsLogTopic(c *gin.Context) {
	var param request.AddMerchantApmRequest
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "创建成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := t.service.AddMerchantClsLogTopic(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

// DeleteMerchantLog godoc
// @Summary 删除商户日志主题和日志集
// @Security ApiKeyAuth
// @Description 删除商户日志主题和日志集
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteMerchantLog true "商户名称，商户id"
// @Tags ThirdTencent
// @Success 200 {object}  response.Response
// @Router /third/tencent/cls/topic/deleteMerchantLog [post]
func (t *ApiCls) DeleteMerchantLog(c *gin.Context) {
	var param request.DeleteMerchantLog
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "删除成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := t.service.DeleteMerchantLog(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

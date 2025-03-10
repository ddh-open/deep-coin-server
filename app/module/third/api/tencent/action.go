package tencent

import (
	"devops-http/app/contract"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
)

// TencentListResource godoc
// @Summary 获得腾讯云资源列表接口
// @Security ApiKeyAuth
// @Description 获得腾讯云资源列表接口
// @accept application/json
// @Produce application/json
// @Param data body request.TencentResourceListRequest true "页数，页大小，筛选条件"
// @Tags ThirdTencent
// @Success 200 {object}  response.Response
// @Router /third/tencent/resource/list [post]
func (t *ApiTencent) TencentListResource(c *gin.Context) {
	var param request.TencentResourceListRequest
	err := c.ShouldBindJSON(&param)
	res := response.Response{Code: 1, Msg: "查询成功", Data: nil}
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	result, err := t.service.GetTencentResourceList(param, c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc))
	if err != nil {
		res.Msg = err.Error()
		c.DJson(res)
		return
	}
	res.Data = result
	c.DJson(res)
}

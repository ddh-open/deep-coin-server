package identity

import (
	"devops-http/app/module/base/response"
	"devops-http/app/module/workflow/service/identity"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type apiIdentity struct {
	service *identity.Service
}

func NewApiIdentity(c framework.Container) *apiIdentity {
	return &apiIdentity{service: identity.NewService(c)}
}

// GetParticipant godoc
// @Summary 获取流程的所有参与者
// @Security ApiKeyAuth
// @Description 获取流程的所有参与者
// @accept application/json
// @Produce application/json
// @Param id path int true "流程id"
// @Tags process
// @Success 200 {object}  response.Response
// @Router /v1/workflow/identity/{id} [post]
func (api *apiIdentity) GetParticipant(c *gin.Context) {
	logger := c.MustMakeLog()
	id := c.Param("id")
	res := response.Response{
		Code: 1,
		Msg:  "",
		Data: nil,
	}
	data, err := api.service.GetParticipant(id)
	if err != nil {
		logger.Error(err.Error())
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = data
	}
	c.DJson(res)
}

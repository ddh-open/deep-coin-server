package tencent

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/third"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/resources/proto/thirdGrpc"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	repository *base.Repository
}

func NewService(c framework.Container) *Service {
	db, err := c.MustMake(contract2.ORMKey).(contract2.ORMService).GetDB()
	logger := c.MustMake(contract2.LogKey).(contract2.Log)
	if err != nil {
		logger.Error("service 获取db出错： err", zap.Error(err))
	}
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetRepository() *base.Repository {
	return s.repository
}

func (s *Service) SetRepository(model interface{}) *base.Repository {
	return s.repository.SetRepository(model)
}

func (s *Service) GetTencentResourceList(request third.TencentResourceListRequest, grpcService contract.ServiceGrpc, param ...interface{}) (response.PageResult, error) {
	conn, err := grpcService.GetGrpc("grpc.third")
	var result response.PageResult
	var list []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := thirdGrpc.NewTencentServiceClient(conn)
	resp, err := client.GetResourceList(context.Background(), &thirdGrpc.ThirdTencentResourceRequest{
		Type:   request.Type,
		Ids:    request.Ids,
		Names:  request.Names,
		Offset: request.Offset,
		Limit:  request.Limit,
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() == 200 {
		err = json.Unmarshal(resp.GetList(), &list)
	} else {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}
	result.List = list
	result.PageSize = resp.GetPageSize()
	result.Page = resp.GetPage()
	result.Total = resp.GetCounts()
	return result, err
}

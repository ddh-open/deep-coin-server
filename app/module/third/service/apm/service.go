package apm

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/resources/proto/thirdGrpc"
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

func (s *Service) AddMerchantApm(request request.AddMerchantApmRequest, grpcService contract.ServiceGrpc, param ...interface{}) (result response.Response, err error) {
	conn, err := grpcService.GetGrpc("grpc.third")
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := thirdGrpc.NewApmServiceClient(conn)
	type res struct {
		Token                   string `json:"token"`
		InstanceId              string `json:"instance_id"`
		PrivateLinkCollectorURL string `json:"private_link_collector_url"`
		Namespace               string `json:"namespace"`
	}
	var data []*res
	for _, namespace := range request.Namespaces {
		resp, err := client.CreateMerchantApmInstance(context.Background(), &thirdGrpc.CreateMerchantApmInstanceRequest{
			MerchantId:   request.MerchantId,
			Namespace:    request.MerchantName,
			MerchantName: namespace,
		})
		if err != nil {
			return result, err
		}
		data = append(data, &res{
			Token:                   resp.Token,
			InstanceId:              resp.InstanceId,
			PrivateLinkCollectorURL: resp.PrivateLinkCollectorURL,
			Namespace:               namespace,
		})
		result.Data = data
	}
	return result, nil
}

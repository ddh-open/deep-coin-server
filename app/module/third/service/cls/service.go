package cls

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/third/model/cls"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/resources/proto/thirdGrpc"
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

func (s *Service) DeleteMerchantLog(request request.DeleteMerchantLog, grpcService contract.ServiceGrpc, param ...interface{}) (result response.Response, err error) {
	// 1. 创建日志集，得到日志集id
	conn, err := grpcService.GetGrpc("grpc.third")
	defer conn.Close()
	client := thirdGrpc.NewClsServiceClient(conn)
	if err != nil {
		return result, err
	}

	resp, err := client.DeleteMerchantLog(context.Background(), &thirdGrpc.DeleteMerchantLogRequest{
		MerchantName: request.MerchantName,
		MerchantId:   request.MerchantId,
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() != 200 {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}

	result.Msg = resp.GetResult().Msg

	return result, err

}

func (s *Service) AddMerchantClsLogTopic(request request.AddMerchantApmRequest, grpcService contract.ServiceGrpc, param ...interface{}) (result response.Response, err error) {
	// 1. 创建日志集，得到日志集id
	conn, err := grpcService.GetGrpc("grpc.third")
	defer conn.Close()
	client := thirdGrpc.NewClsServiceClient(conn)
	if err != nil {
		return result, err
	}
	resp, err := client.CreateLogset(context.Background(), &thirdGrpc.LogsetCreateRequest{
		LogsetName: request.MerchantName, // LogsetName == MerchantName
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() != 200 {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}

	logsetId := resp.LogsetId // 获取商户的日志id

	// 2. 创建每个名称空间的日志主题，得到日志主题信息
	res, err := client.BatchCreateClsTopic(context.Background(), &thirdGrpc.BatchCreateClsTopicRequest{
		LogsetId:     logsetId,
		TopicName:    request.Namespaces,
		MerchantId:   request.MerchantId,
		MerchantName: request.MerchantName,
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() != 200 {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}

	var topics []*cls.Topic
	for _, v := range res.GetTopics() {
		topics = append(topics, &cls.Topic{
			LogsetId:  v.GetLogsetId(),
			TopicId:   v.GetTopicId(),
			TopicName: v.GetTopicName(),
		})
	}
	result.Msg = resp.Result.GetMsg()
	result.Data = response.AddMerchantClsLogTopicRequest{
		MerchantName: request.MerchantName,
		MerchantId:   request.MerchantId,
	}

	return result, nil
}

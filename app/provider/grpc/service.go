package grpc

import (
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	container framework.Container
}

func NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &Service{container: container}, nil
}

func (s *Service) GetGrpc(configPath string, opt ...interface{}) (*grpc.ClientConn, error) {
	config := s.container.MustMake(contract2.ConfigKey).(contract2.Config)
	for _, v := range opt {
		if token, ok := v.(credentials.PerRPCCredentials); ok {
			return grpc.Dial(config.GetString(configPath+".port"), grpc.WithPerRPCCredentials(token))
		}
	}
	return grpc.Dial(config.GetString(configPath+".port"), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

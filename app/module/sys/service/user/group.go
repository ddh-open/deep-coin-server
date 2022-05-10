package user

import (
	"context"
	"devops-http/app/contract"
	"devops-http/resources/proto/userGrpc"
	"encoding/json"
	"github.com/pkg/errors"
)

func (s *Service) GetGroupsByUserId(id string, domain string, grpcService contract.ServiceGrpc, param ...interface{}) ([]map[string]interface{}, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceCabinClient(conn)
	resp, err := client.GetCabinById(context.Background(), &userGrpc.WithSliderParamRequest{
		PType:      "g2",
		FieldIndex: 0,
		FieldValue: []string{id, "", domain},
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() == 200 {
		err = json.Unmarshal(resp.GetData(), &result)
	} else {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}
	return result, err
}

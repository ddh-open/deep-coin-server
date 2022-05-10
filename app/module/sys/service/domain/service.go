package domain

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/resources/proto/userGrpc"
	"encoding/json"
	"github.com/ddh-open/gin/framework"
	contract2 "github.com/ddh-open/gin/framework/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
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

func (s *Service) GetDomainById(id string, grpcService contract.ServiceGrpc) ([]map[string]interface{}, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceDomainClient(conn)
	resp, err := client.DomainList(context.Background(), &userGrpc.ListRequest{
		Filter: []string{"id = ?", id},
	})
	if err != nil {
		return result, err
	}
	if resp.GetResult().GetCode() == 200 {
		err = json.Unmarshal(resp.GetList(), &result)
	} else {
		err = errors.Wrap(err, resp.GetResult().GetMsg())
	}
	return result, err
}

func (s *Service) GetDomainList(request request.PageRequest, grpcService contract.ServiceGrpc, param ...interface{}) (response.PageResult, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result response.PageResult
	var list []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceDomainClient(conn)
	var filter []string
	if len(request.Filter) > 0 {
		filter = append(filter, "name like ? or english_name like ? or domain_num like ?")
		filter = append(filter, "%"+request.Filter[0]+"%")
		filter = append(filter, "%"+request.Filter[0]+"%")
		filter = append(filter, "%"+request.Filter[0]+"%")
	}
	resp, err := client.DomainList(context.Background(), &userGrpc.ListRequest{
		Filter:   filter,
		Page:     request.Page,
		PageSize: request.PageSize,
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
	result.Page = resp.GetPageSize()
	result.Total = resp.GetCounts()
	return result, err
}

func (s *Service) AddDomain(mapData map[string]interface{}, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(&mapData)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceDomainClient(conn)
	resp, err := client.DomainAdd(context.Background(), &userGrpc.BytesRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}

func (s *Service) ModifyDomain(mapData map[string]interface{}, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(&mapData)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceDomainClient(conn)
	resp, err := client.DomainModify(context.Background(), &userGrpc.BytesRequest{
		Data: data,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}

func (s *Service) DeleteDomain(ids string, grpcService contract.ServiceGrpc, param ...interface{}) error {
	var idsInt []int64
	if strings.Contains(ids, ",") {
		for _, s2 := range strings.Split(ids, ",") {
			idsInt = append(idsInt, cast.ToInt64(s2))
		}
	} else {
		idsInt = append(idsInt, cast.ToInt64(ids))
	}
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	client := userGrpc.NewServiceDomainClient(conn)
	resp, err := client.DomainDelete(context.Background(), &userGrpc.IdsRequest{
		Ids: idsInt,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}

package role

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/role"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/resources/proto/userGrpc"
	"encoding/json"
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

func (s *Service) GetRolesResource(name string, domain string, grpcService contract.ServiceGrpc, param ...interface{}) ([]map[string]interface{}, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceCabinClient(conn)
	resp, err := client.GetCabinById(context.Background(), &userGrpc.WithSliderParamRequest{
		PType:      "p",
		FieldIndex: 0,
		FieldValue: []string{name, domain, "", ""},
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

func (s *Service) GetRoleById(id string, grpcService contract.ServiceGrpc) ([]map[string]interface{}, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceRoleClient(conn)
	resp, err := client.RoleList(context.Background(), &userGrpc.ListRequest{
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

func (s *Service) GetRoleList(request request.PageRequest, grpcService contract.ServiceGrpc, param ...interface{}) (response.PageResult, error) {
	conn, err := grpcService.GetGrpc("grpc.user")
	var result response.PageResult
	var list []map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewServiceRoleClient(conn)
	resp, err := client.RoleList(context.Background(), &userGrpc.ListRequest{
		Filter:   request.Filter,
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

func (s *Service) AddRole(dataMap map[string]interface{}, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(&dataMap)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceRoleClient(conn)
	resp, err := client.RoleAdd(context.Background(), &userGrpc.BytesRequest{
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

func (s *Service) ModifyRole(mapData map[string]interface{}, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(&mapData)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceRoleClient(conn)
	resp, err := client.RoleModify(context.Background(), &userGrpc.BytesRequest{
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

func (s *Service) DeleteRole(ids string, grpcService contract.ServiceGrpc, param ...interface{}) error {
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
	client := userGrpc.NewServiceRoleClient(conn)
	resp, err := client.RoleDelete(context.Background(), &userGrpc.IdsRequest{
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

func (s *Service) AddResourcesToRole(request []request.CabinInReceive, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	requestBytes, err := json.Marshal(&request)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceCabinClient(conn)
	resp, err := client.CabinRuleAdd(context.Background(), &userGrpc.BytesRequest{
		Data: requestBytes,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}

// CopyRole 复制角色
func (s *Service) CopyRole(req request.CopyRoleParams, c contract.Cabin) error {
	var copyRole role.DevopsSysRole
	s.repository.GetDB().Model(&role.DevopsSysRole{}).Where("id = ?", req.CopyId).First(&copyRole)
	if copyRole.ID <= 0 {
		return errors.New("未找到要复制的角色")
	}
	err := s.repository.GetDB().Create(&req.DevopsSysRole).Error
	if err != nil {
		return err
	}
	list := c.GetCabin().GetFilteredPolicy(0, req.CopyId)
	for _, i := range list {
		i[0] = cast.ToString(req.ID)
	}
	_, err = c.GetCabin().AddPolicies(list)
	return err
}

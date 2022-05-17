package role

import (
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/role"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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

func (s *Service) GetRoleById(id string) (result role.DevopsSysRoleView, err error) {
	var data role.DevopsSysRole
	err = s.repository.GetDB().Model(&role.DevopsSysRole{}).Where("id = ?", id).First(&data).Error
	result.DevopsSysRole = data
	return
}

func (s *Service) GetRoleList(req request.SearchRoleParams, userToken *base.TokenUser, cabin contract.Cabin) (result response.PageResult, err error) {
	lists := make([]role.DevopsSysRole, 0)
	db := s.repository.GetDB().Model(&role.DevopsSysRole{})
	if req.Name != "" {
		db.Where(" name like ? ", "%"+req.Name+"%")
	}
	if req.ID != 0 {
		db.Where(" id = ? ", req.ID)
	}
	err = db.Count(&result.Total).Error
	if err != nil {
		return result, err
	}
	result.PageSize = req.PageSize
	result.Page = req.Page
	err = db.Limit(int(req.PageSize)).Offset(int((req.Page - 1) * req.PageSize)).Order("id desc").Find(&lists).Error
	if err != nil {
		return result, err
	}
	resultList := make([]role.DevopsSysRoleView, 0)
	for _, list := range lists {
		var handleRoleData role.DevopsSysRoleView
		handleRoleData.DevopsSysRole = list
		cabin.GetCabin().ClearPolicy()
		data := cabin.GetCabin().GetFilteredNamedPolicy("p", 0, cast.ToString(list.ID), userToken.CurrentDomain, "", "")
		for _, datum := range data {
			if datum[3] == base.SourceList[base.MENUS] {
				handleRoleData.Menus = append(handleRoleData.Menus, cast.ToInt(datum[2]))
			}
			if datum[3] == base.SourceList[base.APIS] {
				handleRoleData.Apis = append(handleRoleData.Apis, cast.ToInt(datum[2]))
			}
		}
		resultList = append(resultList, handleRoleData)
	}
	result.List = resultList
	return result, err
}

func (s *Service) GetRoleTree() (result response.PageResult, err error) {
	lists := make([]role.DevopsSysRole, 0)
	db := s.repository.GetDB().Model(&role.DevopsSysRole{})
	err = db.Count(&result.Total).Error
	if err != nil {
		return result, err
	}
	err = db.Order("sort desc").Find(&lists).Error
	if err != nil {
		return result, err
	}
	result.List = lists
	return result, err
}

func (s *Service) AddRole(req role.DevopsSysRoleEntity) error {
	return s.repository.GetDB().Model(role.DevopsSysRole{}).Create(&req.DevopsSysRole).Error
}

func (s *Service) ModifyRole(req role.DevopsSysRoleEntity) error {
	return s.repository.GetDB().Model(role.DevopsSysRole{}).Where("id = ?", req.ID).Save(&req.DevopsSysRole).Error
}

func (s *Service) DeleteRole(ids string) error {
	var roles []role.DevopsSysRole
	err := s.repository.GetDB().Model(&role.DevopsSysRole{}).Where("id in (?)", ids).Find(&roles).Error
	if err != nil {
		return err
	}
	for _, sysRole := range roles {
		err = s.repository.GetDB().Model(&role.DevopsSysRole{}).Unscoped().Delete(&sysRole).Error
		if err != nil {
			return err
		}
	}
	return nil
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

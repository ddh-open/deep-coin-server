package user

import (
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/sys"
	"devops-http/app/module/sys/model/path"
	"devops-http/app/module/sys/model/role"
	"devops-http/app/module/sys/model/user"
	"github.com/pkg/errors"
)

func (s *Service) GetRolesByUserId(id string, domain string, c contract.Cabin) ([]role.DevopsSysRole, error) {
	result := make([]role.DevopsSysRole, 0)
	var userData user.DevopsSysUser
	s.repository.SetRepository(&user.DevopsSysUser{}).First(&userData, "id = "+id)
	if userData.ID <= 0 {
		err := errors.New("未找到该用户！")
		return result, err
	}
	roleIds, err := c.GetCabin().GetRolesForUser(userData.UUID.String(), domain)
	if err != nil {
		return result, err
	}
	err = s.repository.SetRepository(&role.DevopsSysRole{}).GetDB().Find(&result, "id in (?)", roleIds).Error
	if err != nil {
		return result, err
	}
	return result, err
}

func (s *Service) GetUserApis(userToken *base.TokenUser, c contract.Cabin) (result []path.DevopsSysApi, err error) {
	list, err := c.GetCabin().GetImplicitResourcesForUser(userToken.Uuid, userToken.CurrentDomain)
	if err != nil {
		return result, err
	}
	var apiIds []string
	for _, str := range list {
		if str[3] == "APIS" {
			apiIds = append(apiIds, str[2])
		}
	}
	err = s.repository.SetRepository(&path.DevopsSysApi{}).GetDB().Find(&result, "id in (?)", apiIds).Error
	if err != nil {
		return result, err
	}
	return result, err
}

func (s *Service) RelativeRolesToUser(request sys.RelativeUserRequest, domain string, c contract.Cabin) error {
	userData := user.DevopsSysUser{}
	s.repository.GetDB().First(&userData, "id = ?", request.UserId)
	if userData.ID <= 0 {
		return errors.New("未找到该用户！")
	}
	_, err := c.GetCabin().DeleteRolesForUser(userData.UUID.String(), domain)
	if err != nil {
		return err
	}
	_, err = c.GetCabin().AddRolesForUser(userData.UUID.String(), request.RoleIds, domain)
	if err != nil {
		return errors.New("关联失败！")
	}
	return err
}

func (s *Service) DeleteRelativeRolesToUser(request sys.RelativeUserRequest, domain string, c contract.Cabin) error {
	userData := user.DevopsSysUser{}
	s.repository.GetDB().First(&userData, "id = ?", request.UserId)
	if userData.ID <= 0 {
		return errors.New("未找到该用户！")
	}
	for i := range request.RoleIds {
		_, err := c.GetCabin().DeleteRoleForUser(userData.UUID.String(), request.RoleIds[i], domain)
		if err != nil {
			return errors.New("删除失败！")
		}
	}
	return nil
}

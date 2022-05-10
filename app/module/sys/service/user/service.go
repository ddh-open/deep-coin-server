package user

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/sys"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/sys/model/config"
	"devops-http/app/module/sys/model/role"
	"devops-http/app/module/sys/model/user"
	"devops-http/resources/proto/userGrpc"
	"encoding/base64"
	"fmt"
	"github.com/ddh-open/gin/framework"
	contract2 "github.com/ddh-open/gin/framework/contract"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	repository *base.Repository
}

func NewService(c framework.Container) *Service {
	db, err := c.MustMake(contract2.ORMKey).(contract2.ORMService).GetDB()
	logger := c.MustMake(contract2.LogKey).(contract2.Log)
	if err != nil {
		logger.Error("Service 获取db出错： err", zap.Error(err))
	}
	db.AutoMigrate(&user.DevopsSysUser{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetRepository() *base.Repository {
	return s.repository
}

func (s *Service) GetUsers() {

}

func (s *Service) Login(request sys.LoginRequest, grpc contract.ServiceGrpc) (interface{}, error) {
	result := make(map[string]string, 1)
	conn, err := grpc.GetGrpc("grpc.user")
	if err != nil {
		err = errors.Wrap(err, "初始化grpc连接出错")
		return result, err
	}
	defer conn.Close()
	client := userGrpc.NewUserServiceClient(conn)
	resp, err := client.Login(context.Background(), &userGrpc.WithPasswordRequest{
		Username: request.Username,
		Password: request.Password,
		Type:     request.Type,
	})
	if err != nil {
		err = errors.Wrap(err, "grpc 登录接口出错")
		return result, err
	}
	// 代表响应成功
	if resp.GetResult().Code != 200 {
		err = errors.Wrap(errors.New("grpc code -1"), resp.GetResult().GetMsg())
	}
	result["token"] = resp.GetToken()
	return result, err
}

func (s *Service) Modify(req user.DevopsSysUserEntity, l contract.Ldap, c contract.Cabin) (interface{}, error) {
	var oldUser user.DevopsSysUser
	err := s.repository.SetRepository(&user.DevopsSysUser{}).Find(&oldUser, "id = ?", req.ID)
	if err != nil {
		return nil, errors.Errorf("未找到需要编辑的用户：%s", err.Error())
	}
	if oldUser.UserType == 1 && req.Password != "" {
		err = l.ChangePassword(req.WorkNum, req.Password)
		if err != nil {
			return nil, errors.Errorf("密码修改失败: %s", err)
		}
		req.Password = ""
	} else if req.Password != "" {
		passwd, err := base64.StdEncoding.DecodeString(req.Password)
		if err != nil {
			return req.DevopsSysUser, err
		}
		req.DevopsSysUser.Password = utils.MD5V(passwd)
	}
	err = s.repository.SetRepository(&user.DevopsSysUser{}).Update(&req.DevopsSysUser, "id = ?", req.ID)
	req.Password = ""
	// 删除之前的角色
	flag, err := c.GetCabin().DeleteRolesForUser(oldUser.UUID.String(), oldUser.Merchants)
	if !flag {
		err = errors.New("删除角色失败")
		return req.DevopsSysUser, err
	}
	// 添加角色
	flag, err = c.GetCabin().AddRolesForUser(oldUser.UUID.String(), req.RoleIds, oldUser.Merchants)
	if !flag {
		err = errors.New("增加角色失败")
	}
	err = s.repository.SetRepository(&user.DevopsSysUser{}).Find(&oldUser, "id = ?", req.ID)
	return oldUser, err
}

func (s *Service) Add(req user.DevopsSysUserEntity, l contract.Ldap, c contract.Cabin) (interface{}, error) {
	userData := req.DevopsSysUser
	passwd, err := base64.StdEncoding.DecodeString(userData.Password)
	if err != nil {
		return userData, err
	}
	userData.Password = utils.MD5V(passwd)
	userData.UUID = uuid.NewV4()
	if !errors.Is(s.repository.GetDB().Where("name = ? ", req.Username).First(&config.DevopsSysConfig{}).Error, gorm.ErrRecordNotFound) {
		return userData, errors.New("存在相同用户名的用户")
	}
	if userData.UserType == 1 {
		filter := "OU=" + req.Merchants
		if req.Merchants != "freemud" {
			filter += ",OU=Merchants"
		}
		// ad 账户
		err = l.CreateUser(userData.Username, fmt.Sprintf("%v,DC=office,DC=freemud,DC=cn", filter), string(passwd), nil)
		if err != nil {
			return userData, errors.Errorf("AD 账户新增失败：%s", err.Error())
		}
	}
	err = s.repository.SetRepository(&user.DevopsSysUser{}).Save(&userData)
	if err != nil {
		return nil, errors.Errorf("新增失败：%s", err.Error())
	}
	userData.Password = ""
	// 添加角色
	flag, err := c.GetCabin().AddRolesForUser(userData.UUID.String(), req.RoleIds, userData.Merchants)
	if !flag {
		err = errors.New("增加角色失败")
	}
	return userData, err
}

func (s *Service) Delete(ids string, l contract.Ldap, c contract.Cabin) error {
	var users []user.DevopsSysUser
	err := s.repository.SetRepository(&user.DevopsSysUser{}).Find(&users, "id in (?)", ids)
	if err != nil {
		return errors.Errorf("未找到需要删除的用户：%s", err.Error())
	}
	if len(users) <= 0 {
		return errors.Errorf("未找到需要删除的用户")
	}
	for _, sysUser := range users {
		if sysUser.UserType == 1 {
			filter := "OU=" + sysUser.Merchants
			if sysUser.Merchants != "freemud" {
				filter += ",OU=Merchants"
			}
			err = l.DeleteUser(sysUser.Username, fmt.Sprintf("%v,DC=office,DC=freemud,DC=cn", filter))
			if err != nil {
				return errors.Errorf("AD删除用户: %s出错：%s", sysUser.Username, err.Error())
			}
		}
		err = s.repository.SetRepository(&user.DevopsSysUser{}).GetDB().Where("id = ?", sysUser.ID).Delete(&sysUser).Error
		if err != nil {
			return errors.Errorf("数据库删除用户: %s出错：%s", sysUser.Username, err.Error())
		}
		// 删除之前的角色
		flag, err := c.GetCabin().DeleteRolesForUser(sysUser.UUID.String(), sysUser.Merchants)
		if !flag {
			err = errors.New("删除角色失败")
			return err
		}
	}
	return err
}

func (s *Service) ChangePassword(req sys.ChangePasswordRequest, l contract.Ldap) (err error) {
	if req.Type == 1 {
		_, err = l.Login(req.Username, req.OldPassword)
		if err != nil {
			return errors.Errorf("原密码不正确: %s", err)
		}
		err = l.ChangePassword(req.Username, req.Password)
		if err != nil {
			return errors.Errorf("密码修改失败: %s", err)
		}
	}
	return
}

// UserList 获取用户列表
func (s *Service) UserList(e contract.Cabin, req request.PageRequest) (interface{}, error) {
	res := make([]user.DevopsSysUser, 0)
	result := response.PageResult{
		List:     nil,
		Columns:  nil,
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	err := s.repository.SetRepository(&user.DevopsSysUser{}).List(&res, req.PageSize, req.Page, req.Filter)
	if err != nil {
		return nil, err
	}
	err = s.repository.SetRepository(&user.DevopsSysUser{}).Counts(&result.Total, req.Filter)
	if err != nil {
		return nil, err
	}
	resView := make([]user.DevopsSysUserView, 0)
	for _, re := range res {
		roleList, _ := e.GetCabin().GetRolesForUser(re.UUID.String(), re.Merchants)
		resView = append(resView, user.DevopsSysUserView{DevopsSysUser: re, RoleIds: roleList})
	}
	result.List = resView
	result.Columns = user.SysUserViewColumns
	return result, err
}

// UserInfo 获取用户详细信息
func (s *Service) UserInfo(token *base.TokenUser, e contract.Cabin, filter []interface{}) (user.DevopsSysUserView, error) {
	res := user.DevopsSysUser{}
	err := s.repository.SetRepository(&user.DevopsSysUser{}).GetDB().First(&res, filter...).Error
	resView := user.DevopsSysUserView{DevopsSysUser: res}
	if res.ID <= 0 {
		return resView, errors.Errorf("未找到该用户： %v ！", err)
	}
	roleList, _ := e.GetCabin().GetRolesForUser(resView.UUID.String(), token.CurrentDomain)
	resView.RoleIds = roleList
	var roleData []role.DevopsSysRole
	s.repository.SetRepository(&role.DevopsSysRole{}).Find(&roleData, roleList)
	for i := range roleData {
		if i == len(roleData)-1 {
			resView.RoleName += roleData[i].Name
		} else {
			resView.RoleName += roleData[i].Name + ","
		}
	}
	return resView, err
}

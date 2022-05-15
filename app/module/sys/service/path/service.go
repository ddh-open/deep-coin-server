package path

import (
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/sys/model/path"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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
		logger.Error("service 获取db出错： err", zap.Error(err))
	}
	return &Service{base.NewRepository(db)}
}

func (s *Service) CreateApi(api path.DevopsSysApi) (err error) {
	if !errors.Is(s.repository.GetDB().Where("path = ? AND method = ?", api.Path, api.Method).First(&path.DevopsSysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return s.repository.GetDB().Create(&api).Error
}

func (s *Service) DeleteApi(api path.DevopsSysApi, c contract.Cabin) (err error) {
	err = s.repository.GetDB().Delete(&api).Error
	c.ClearCabin(2, api.Path, api.Method)
	return err
}

func (s *Service) GetAPIInfoList(req request.SearchApiParams, logGet contract2.Log) (err error, list interface{}, total int64) {
	limit := int(req.PageSize)
	offset := int(req.PageSize * (req.Page - 1))
	db := s.repository.GetDB().Model(&path.DevopsSysApi{})
	var apiList []path.DevopsSysApi

	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}

	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}

	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}

	if req.ApiGroup != "" {
		db = db.Where("api_group = ?", req.ApiGroup)
	}

	err = db.Count(&total).Error

	if err != nil {
		return err, apiList, total
	} else {
		db = db.Limit(limit).Offset(offset)
		if req.OrderKey != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["path"] = true
			orderMap["api_group"] = true
			orderMap["description"] = true
			orderMap["method"] = true
			if orderMap[req.OrderKey] {
				if req.Desc {
					OrderStr = req.OrderKey + " desc"
				} else {
					OrderStr = req.OrderKey
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", req.OrderKey)
				return err, apiList, total
			}

			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			err = db.Order("api_group").Find(&apiList).Error
		}
	}
	return err, apiList, total
}

func (s *Service) GetAllApis() (err error, result map[string][]path.DevopsSysApi) {
	var apis []path.DevopsSysApi
	err = s.repository.GetDB().Find(&apis).Error
	result = make(map[string][]path.DevopsSysApi, 0)
	for _, sysApis := range apis {
		if data, ok := result[sysApis.ApiGroup]; ok {
			data = append(data, sysApis)
			result[sysApis.ApiGroup] = data
		} else {
			var list []path.DevopsSysApi
			list = append(list, sysApis)
			result[sysApis.ApiGroup] = list
		}
	}
	return
}

func (s *Service) GetApiById(id string) (err error, api path.DevopsSysApi) {
	err = s.repository.GetDB().Where("id = ?", id).First(&api).Error
	return
}

func (s *Service) RelativeApiToRole(req request.RelativeRoleApisRequest, userToken *base.TokenUser, cabin contract.Cabin) (err error) {
	if len(req.ApiIds) <= 0 {
		return errors.New("api id为空！")
	}
	var apis []path.DevopsSysApi
	s.repository.GetDB().Find(&apis, "id in (?)", req.ApiIds)
	if len(apis) <= 0 {
		return errors.New("api 校验失败！")
	}
	cabin.GetCabin().ClearPolicy()
	_, err = cabin.GetCabin().RemoveFilteredNamedPolicy("p", 0, req.RoleId, userToken.CurrentDomain, "", "APIS")
	if err != nil {
		return err
	}
	var rule [][]string
	for i := range apis {
		rule = append(rule, []string{req.RoleId, userToken.CurrentDomain, cast.ToString(apis[i].ID), "APIS", apis[i].Method})
	}
	// 给角色添加菜单树
	cabin.GetCabin().ClearPolicy()
	_, err = cabin.GetCabin().AddPolicies(rule)
	return err
}

func (s *Service) UpdateApi(api path.DevopsSysApi, c contract.Cabin) (err error) {
	var oldA path.DevopsSysApi
	err = s.repository.GetDB().Where("id = ?", api.ID).First(&oldA).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		if !errors.Is(s.repository.GetDB().Where("path = ? AND method = ?", api.Path, api.Method).First(&path.DevopsSysApi{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同api路径")
		}
	}
	if err != nil {
		return err
	} else {
		err = c.UpdateCabinApi(oldA.Path, api.Path, oldA.Method, api.Method)
		if err != nil {
			return err
		} else {
			err = s.repository.GetDB().Where("id = ?", api.ID).Save(&api).Error
		}
	}
	return err
}

func (s *Service) DeleteApisByIds(ids request.ReqById) (err error) {
	err = s.repository.GetDB().Delete(&[]path.DevopsSysApi{}, "id in (?)", ids.Ids).Error
	return err
}

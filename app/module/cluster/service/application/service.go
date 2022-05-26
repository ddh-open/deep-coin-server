package application

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/cluster/model/application"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
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
	db.AutoMigrate(&application.DevopsClusterApplication{})
	db.AutoMigrate(&application.DevopsClusterApplicationConfig{})
	err = db.SetupJoinTable(&base.DevopsSysGroup{}, "Groups", &application.DevopsClusterApplicationRelativeGroup{})
	err = db.SetupJoinTable(application.DevopsClusterApplicationConfig{}, "Configs", &application.DevopsClusterApplicationRelativeConfig{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) List(req request.SearchApplicationParams) (result response.PageResult, err error) {
	result.Page = req.Page
	result.PageSize = req.PageSize
	limit := int(result.PageSize)
	offset := int(result.PageSize * (result.Page - 1))
	db := s.repository.GetDB().Model(&application.DevopsClusterApplication{})
	var applicationList []application.DevopsClusterApplication
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	err = db.Count(&result.Total).Error

	if err != nil {
		return
	} else {
		db = db.Limit(limit).Offset(offset)
		if req.OrderKey != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 4)
			orderMap["id"] = true
			orderMap["name"] = true
			if orderMap[req.OrderKey] {
				if req.Desc {
					OrderStr = req.OrderKey + " desc"
				} else {
					OrderStr = req.OrderKey
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", req.OrderKey)
				return
			}
			err = db.Order(OrderStr).Find(&applicationList).Error
		} else {
			err = db.Order("id").Find(&applicationList).Error
		}
	}
	for i := range applicationList {
		s.repository.GetDB().Model(&applicationList[i]).Association("Groups").Find(&applicationList[i].Groups)
		s.repository.GetDB().Model(&applicationList[i]).Association("Configs").Find(&applicationList[i].Configs)
	}
	result.List = applicationList
	return result, err
}

func (s *Service) GetApplicationById(id string) (result application.DevopsClusterApplication, err error) {
	err = s.repository.GetDB().Model(&application.DevopsClusterApplication{}).Where("id = ?", id).First(&result).Association("Groups").Find(&result.Groups)
	if err != nil {
		return
	}
	err = s.repository.GetDB().Model(&result).Association("Configs").Find(&result.Configs)
	return
}

func (s *Service) Save(req *application.DevopsClusterApplication) (err error) {
	err = s.repository.GetDB().Model(&application.DevopsClusterApplication{}).Create(req).Error
	return
}

func (s *Service) Modify(req *application.DevopsClusterApplication) (err error) {
	err = s.repository.GetDB().Model(&application.DevopsClusterApplication{}).Where("id = ?", req.ID).Save(req).Error
	return
}

func (s *Service) ModifyConfig(req *application.DevopsClusterApplication) (err error) {
	if len(req.Configs) <= 0 {
		return errors.New("应用环境配置为空")
	}
	err = s.repository.GetDB().Model(&application.DevopsClusterApplicationConfig{}).Where("id = ?", req.Configs[0].ID).Save(&req.Configs[0]).Error
	return
}

func (s *Service) ModifyGroup(req *application.DevopsClusterApplication) (err error) {
	if len(req.Groups) <= 0 {
		return errors.New("应用分组配置为空")
	}
	err = s.repository.GetDB().Model(&req).Association("Groups").Replace(base.DevopsSysGroup{DevopsModel: base.DevopsModel{ID: req.Groups[0].ID}}, req.Groups[0])
	return
}

func (s *Service) AddGroup(req *application.DevopsClusterApplication) (err error) {
	if len(req.Groups) <= 0 {
		return errors.New("应用分组配置为空")
	}
	err = s.repository.GetDB().Model(&req).Association("Groups").Append(req.Groups)
	return
}

func (s *Service) AddConfig(req *application.DevopsClusterApplication) (err error) {
	if len(req.Configs) <= 0 {
		return errors.New("应用环境配置为空")
	}
	err = s.repository.GetDB().Model(&req).Association("Configs").Append(req.Configs)
	return
}

func (s *Service) DeleteGroup(req *application.DevopsClusterApplication) (err error) {
	if len(req.Groups) <= 0 {
		return errors.New("应用分组配置为空")
	}
	err = s.repository.GetDB().Model(&req).Unscoped().Association("Groups").Delete(req.Groups[0])
	return
}

func (s *Service) DeleteConfig(req *application.DevopsClusterApplication) (err error) {
	if len(req.Configs) <= 0 {
		return errors.New("应用环境配置为空")
	}
	err = s.repository.GetDB().Model(&req).Unscoped().Association("Configs").Delete(req.Configs[0])
	return
}

func (s *Service) Delete(id string) (err error) {
	// 删除
	err = s.repository.GetDB().Model(&application.DevopsClusterApplication{}).Unscoped().Delete(&application.DevopsClusterApplication{}, id).Error
	return
}

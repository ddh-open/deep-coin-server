package config

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/sys/model/config"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
	"github.com/pkg/errors"
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
	db.AutoMigrate(config.DevopsSysConfig{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) Create(req config.DevopsSysConfig) (err error) {
	if !errors.Is(s.repository.GetDB().Where("name = ? ", req.Name).First(&config.DevopsSysConfig{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同配置")
	}
	return s.repository.GetDB().Create(&req).Error
}

func (s *Service) Delete(req request.ReqById) (err error) {
	err = s.repository.GetDB().Where("id in (?)", req.Ids).Delete(&config.DevopsSysConfig{}).Error
	return err
}

func (s *Service) List(req request.SearchConfigParams) (err error, list interface{}, total int64) {
	limit := int(req.PageSize)
	offset := int(req.PageSize * (req.Page - 1))
	db := s.repository.GetDB().Model(&config.DevopsSysConfig{})
	var configList []config.DevopsSysConfig

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Key != "" {
		db = db.Where("key LIKE ?", "%"+req.Key+"%")
	}

	if req.Value != "" {
		db = db.Where("value = ?", req.Value)
	}

	if req.Remark != "" {
		db = db.Where("remark = ?", req.Remark)
	}

	err = db.Count(&total).Error

	if err != nil {
		return err, configList, total
	} else {
		db = db.Limit(limit).Offset(offset)
		if req.OrderKey != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["name"] = true
			orderMap["key"] = true
			orderMap["value"] = true
			orderMap["remark"] = true
			if orderMap[req.OrderKey] {
				if req.Desc {
					OrderStr = req.OrderKey + " desc"
				} else {
					OrderStr = req.OrderKey
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", req.OrderKey)
				return err, configList, total
			}

			err = db.Order(OrderStr).Find(&configList).Error
		} else {
			err = db.Order("id").Find(&configList).Error
		}
	}
	return err, configList, total
}

func (s *Service) Update(req config.DevopsSysConfig) (err error) {
	return s.repository.GetDB().Save(&req).Error
}

func (s *Service) GetConfigByName(name string) (data config.DevopsSysConfig, err error) {
	err = s.repository.GetDB().Where("name = ? ", name).First(&data).Error
	return
}

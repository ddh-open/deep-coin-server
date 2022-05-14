package icon

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/icon"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
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

func (s *Service) GetList(req request.SearchIconParams) (list interface{}, err error) {
	var result response.PageResult
	result.Page = req.Page
	result.PageSize = req.PageSize
	limit := int(result.PageSize)
	offset := int(result.PageSize * (result.Page - 1))
	db := s.repository.GetDB().Model(&icon.DevopsSysIcon{})
	var iconList []string
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	err = db.Count(&result.Total).Error
	if err != nil {
		return
	} else {
		err = db.Select("title").Limit(limit).Offset(offset).Order("id").Find(&iconList).Error
	}
	result.List = iconList
	return result, err
}

package operation

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/operation"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
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
	db.AutoMigrate(&operation.DevopsSysOperationRecord{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetDetailById(id string) (data operation.DevopsSysOperationRecord, err error) {
	err = s.repository.GetDB().First(&data, "id = ?", id).Error
	return
}

func (s *Service) GetList(req request.SearchLogsParams) (res response.PageResult, err error) {
	res.Page = req.Page
	res.PageSize = req.PageSize
	limit := int(res.PageSize)
	offset := int(res.PageSize * (res.Page - 1))
	db := s.repository.GetDB().Model(&operation.DevopsSysOperationRecord{})
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Resp != "" {
		db = db.Where("resp LIKE ?", "%"+req.Resp+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", "%"+req.Method+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.Body != "" {
		db = db.Where("body LIKE ?", "%"+req.Body+"%")
	}
	if req.Ip != "" {
		db = db.Where("ip LIKE ?", "%"+req.Ip+"%")
	}
	if req.UserName != "" {
		db = db.Where("username LIKE ?", "%"+req.UserName+"%")
	}

	if len(req.TimeFilter) >= 1 {
		db = db.Where("created_at >= ? and created_at <= ?", req.TimeFilter[0], req.TimeFilter[1])
	}

	err = db.Count(&res.Total).Error
	if err != nil {
		return
	}
	db = db.Limit(limit).Offset(offset)
	var operations []operation.DevopsSysOperationRecord
	if req.OrderKey != "" {
		var OrderStr string
		// 设置有效排序key 防止sql注入
		// 感谢 Tom4t0 提交漏洞信息
		orderMap := make(map[string]bool, 4)
		orderMap["id"] = true
		orderMap["path"] = true
		orderMap["methods"] = true
		orderMap["username"] = true
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
		err = db.Order(OrderStr).Find(&operations).Error
	} else {
		err = db.Order("id desc").Find(&operations).Error
	}
	if err != nil {
		return
	}
	res.List = operations
	return
}

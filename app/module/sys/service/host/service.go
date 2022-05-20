package host

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/sys/model/host"
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
	db.AutoMigrate(&host.DevopsCmdbHost{}, &host.DevopsCmdbHostGroup{})
	err = db.SetupJoinTable(&host.DevopsCmdbHost{}, "Groups", &host.DevopsCmdbHostGroupRelativeHost{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetHostList(req request.SearchHostParams) (result host.DevopsCmdbHostGroup, err error) {
	err = s.repository.GetDB().Model(&host.DevopsCmdbHostGroup{}).Where("id = ?", req.ID).First(&result).Error
	if err != nil {
		return
	}
	err = s.repository.GetDB().Model(&result).Association("Hosts").Find(&result.Hosts)
	return
}

func (s *Service) GetHostGroupTree() (list interface{}, err error) {
	db := s.repository.GetDB().Model(&host.DevopsCmdbHostGroup{})
	db = db.Where("parent_id = ?", 0)
	var groupList []host.DevopsCmdbHostGroup
	err = db.Order("sort").Find(&groupList).Error
	for i := range groupList {
		s.getGroupChildren(&groupList[i])
	}
	return groupList, err
}

func (s *Service) getGroupChildren(group *host.DevopsCmdbHostGroup) {
	var groupList []host.DevopsCmdbHostGroup
	s.repository.GetDB().Model(&host.DevopsCmdbHostGroup{}).Where("parent_id = ?", group.ID).Find(&groupList)
	group.Children = groupList
	for i := range groupList {
		s.getGroupChildren(&groupList[i])
	}
}

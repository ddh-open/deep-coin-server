package identity

import (
	"devops-http/app/module/base"
	identityModel "devops-http/app/module/workflow/model/identity"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
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
	db.AutoMigrate(&identityModel.WorkflowIdentityLink{})
	return &Service{base.NewRepository(db)}
}

func SaveIdentityTx(i *identityModel.WorkflowIdentityLink, tx *gorm.DB) error {
	return tx.Model(&identityModel.WorkflowIdentityLink{}).Save(i).Error
}

func (s *Service) GetParticipant(id string) (result []identityModel.WorkflowIdentityLink, err error) {
	var list []*identityModel.WorkflowIdentityLink
	err = s.repository.GetDB().Model(&identityModel.WorkflowIdentityLink{}).Select("id,user_id,user_name,step,comment").Where("proc_inst_id=? and type=?", id, identityModel.KindsIdentity[identityModel.PARTICIPANT]).Order("id asc").Find(&list).Error
	return
}

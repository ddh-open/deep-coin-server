package task

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/workflow/model/task"
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
	db.AutoMigrate(&task.WorkflowTask{})
	return &Service{base.NewRepository(db)}
}

func SaveTaskByTx(task *task.WorkflowTask, tx *gorm.DB) (uint, error) {
	err := tx.Model(task).Save(task).Error
	return task.ID, err
}

// CompleteTask 审批
func (s *Service) CompleteTask(req *request.TaskReceiver) (bool, error) {
	return true, nil
}

// WithDrawTask 撤回任务
func (s *Service) WithDrawTask(req *request.TaskReceiver) (bool, error) {
	return true, nil
}

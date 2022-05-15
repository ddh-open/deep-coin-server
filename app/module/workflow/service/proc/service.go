package proc

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/workflow/model/node"
	"devops-http/app/module/workflow/model/proc"
	nodeSvc "devops-http/app/module/workflow/service/node"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type Service struct {
	repository *base.Repository
	saveLock   sync.Mutex
}

func NewService(c framework.Container) *Service {
	db, err := c.MustMake(contract2.ORMKey).(contract2.ORMService).GetDB()
	logger := c.MustMake(contract2.LogKey).(contract2.Log)
	if err != nil {
		logger.Error("service 获取db出错： err", zap.Error(err))
	}
	err = db.AutoMigrate(&proc.WorkflowProc{}, &node.WorkflowNodeParams{})
	return &Service{base.NewRepository(db), sync.Mutex{}}
}

func (s *Service) SaveProc(req *request.ProcRequest) (id uint, err error) {
	if len(req.Name) == 0 {
		err = errors.New("流程名不能为空")
		return
	}
	if req.Resource == nil || len(req.Resource.Name) == 0 {
		err = errors.New("字段 resource 不能为空")
		return
	}
	nodeTypes, err := nodeSvc.GetNodeType(s.repository)
	if err != nil {
		return
	}
	err = s.CheckProcessConfig(req.Resource, nodeTypes...)
	if err != nil {
		return
	}

	// 后面开始进行迁移
	s.saveLock.Lock()
	defer s.saveLock.Unlock()
	oldProc := make([]proc.WorkflowProc, 0)
	err = s.repository.SetRepository(&proc.WorkflowProc{}).Find(&oldProc, " name = ? and company = ? ", req.Name, req.Company)
	if err != nil {
		return
	}
	resource, err := json.Marshal(req.Resource)
	if err != nil {
		return
	}
	p := proc.WorkflowProc{
		Name:     req.Name,
		Version:  0,
		Resource: string(resource),
		UserId:   req.UserId,
		Username: req.Username,
		Company:  req.Company,
	}
	if len(oldProc) == 0 {
		p.Version = 1
		err = s.repository.Save(&p)
		id = p.ID
		return
	}
	tx := s.repository.GetDB().Begin()
	// 保存新版本
	p.Version = oldProc[0].Version + 1
	err = tx.Model(&proc.WorkflowProc{}).Save(&p).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// 转移旧版本
	err = s.MoveProcToHistoryByIDTx(oldProc[0].ID, tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return p.ID, nil
}

func (s *Service) MoveProcToHistoryByIDTx(id uint, tx *gorm.DB) error {
	return tx.Delete(&proc.WorkflowProc{}, "id = ?", id).Error
}

// CheckConditionNode 检查条件节点
func (s *Service) CheckConditionNode(nodes []*node.Node, nodeTypes ...node.WorkflowNodeType) error {
	for _, v := range nodes {
		if v.Properties == nil {
			return errors.New("节点【" + v.NodeID + "】的Properties对象为空值！！")
		}
		if len(v.Properties.Conditions) == 0 {
			return errors.New("节点【" + v.NodeID + "】的Conditions对象为空值！！")
		}
		err := s.CheckProcessConfig(v, nodeTypes...)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckProcessConfig 检查流程配置是否有效
func (s *Service) CheckProcessConfig(node *node.Node, nodeTypes ...node.WorkflowNodeType) error {
	// 节点名称是否有效
	if len(node.NodeID) == 0 {
		return errors.New("节点的【nodeId】不能为空！！")
	}
	// 检查类型是否有效
	if len(node.Type) == 0 {
		return errors.New("节点【" + node.NodeID + "】的类型【type】不能为空")
	}
	var flag = false
	for _, val := range nodeTypes {
		if val.Type == node.Type {
			flag = true
			switch val.RelativeKind {
			case 1:

			case 2:
				// 审批和抄送的节点
				if node.Properties == nil || node.Properties.ActionRules == nil {
					return errors.New("节点【" + node.NodeID + "】的Properties属性不能为空，如：`\"properties\": {\"ActionRules\": [{\"type\": \"target_label\",\"labelNames\": \"人事\",\"memberCount\": 1,\"actType\": \"and\"}],}`")
				}
			case 3:
				// 路由节点
				if node.ConditionNodes != nil { // 存在条件节点
					if len(node.ConditionNodes) == 1 {
						return errors.New("节点【" + node.NodeID + "】条件节点下的节点数必须大于1")
					}
					// 根据条件变量选择节点索引
					err := s.CheckConditionNode(node.ConditionNodes, nodeTypes...)
					if err != nil {
						return err
					}
				} else {
					return errors.New("路由节点应该存在 ConditionNodes 属性")
				}
			case 4:
				// exec 执行
			case 5:
				// 条件判断

			}
			break
		}
	}
	if !flag {
		return errors.New("节点【" + node.NodeID + "】的类型为【" + node.Type + "】，为无效类型")
	}
	// 子节点是否存在
	if node.ChildNode != nil {
		return s.CheckProcessConfig(node.ChildNode, nodeTypes...)
	}
	return nil
}

// ListProc 查询proc 列表
func (s *Service) ListProc(request *request.ProcPageReceiver) (result response.PageResult, err error) {
	var list []proc.WorkflowProc
	var filter []interface{}
	for _, s2 := range request.Filter {
		filter = append(filter, s2)
	}
	result.Page = request.Page
	result.PageSize = request.PageSize
	count := int64(0)
	err = s.repository.SetRepository(&proc.WorkflowProc{}).Counts(&count, filter...)
	if err != nil {
		return
	}
	result.Total = count
	err = s.repository.SetRepository(&proc.WorkflowProc{}).List(&list, request.PageSize, request.Page, filter...)
	result.List = list
	return
}

// DeleteProc 删除proc
func (s *Service) DeleteProc(id string) (err error) {
	err = s.repository.Delete(&proc.WorkflowProc{}, " id = ? ", id)
	return
}

package node

import (
	"container/list"
	"devops-http/app/module/base"
	"devops-http/app/module/workflow/model/execution"
	"devops-http/app/module/workflow/model/node"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
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
	db.AutoMigrate(&node.WorkflowNodeType{})
	db.AutoMigrate(&execution.WorkflowExecution{})
	service := &Service{base.NewRepository(db)}
	err = service.GetNodeType()
	if err != nil {
		logger.Error("service 获取nodeType出错： err", zap.Error(err))
	}
	return service
}

func (s *Service) GetNodeType() error {
	err := s.repository.GetDB().Model(node.WorkflowNodeType{}).Where(" kind = ? ", 1).Find(&node.WorkflowNodeTypesList).Error
	if err != nil {
		return err
	}
	return s.repository.GetDB().Model(node.WorkflowNodeType{}).Where(" kind = ? ", 0).Find(&node.WorkflowNodeKindsList).Error
}

func GetNodeType(repository *base.Repository) ([]node.WorkflowNodeType, error) {
	result := make([]node.WorkflowNodeType, 0)
	err := repository.GetDB().Model(&node.WorkflowNodeType{}).Where(" kind = ? ", 1).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, err
}

// ParseProcessConfig 解析流程定义json数据
func ParseProcessConfig(node *node.Node, variable *map[string]string) (*list.List, error) {
	// defer fmt.Println("----------解析结束--------")
	listData := list.New()
	err := parseProcessConfig(node, variable, listData)
	return listData, err
}

func parseProcessConfig(node *node.Node, variable *map[string]string, list *list.List) (err error) {
	// fmt.Printf("nodeId=%s\n", node.NodeID)
	node.AddToExecutionList(list)
	// 存在条件节点
	if node.ConditionNodes != nil {
		// 如果条件节点只有一个或者条件只有一个，直接返回第一个
		if variable == nil || len(node.ConditionNodes) == 1 {
			err = parseProcessConfig(node.ConditionNodes[0].ChildNode, variable, list)
			if err != nil {
				return err
			}
		} else {
			// 根据条件变量选择节点索引
			condNode, err := GetConditionNode(node.ConditionNodes, variable)
			if err != nil {
				return err
			}
			if condNode == nil {
				//str, _ := util.ToJSONStr(variable)
				return errors.New("节点【" + node.NodeID + "】找不到符合条件的子节点,检查变量【var】值是否匹配,")
				// panic(err)
			}
			err = parseProcessConfig(condNode, variable, list)
			if err != nil {
				return err
			}

		}
	}
	// 存在子节点
	if node.ChildNode != nil {
		err = parseProcessConfig(node.ChildNode, variable, list)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetConditionNode 获取条件节点
func GetConditionNode(nodes []*node.Node, maps *map[string]string) (result *node.Node, err error) {
	map2 := *maps
	for _, n := range nodes {
		var flag int
		for _, v := range n.Properties.Conditions[0] {
			paramValue := map2[v.ParamKey]
			if len(paramValue) == 0 {
				return nil, errors.New("流程启动变量【var】的key【" + v.ParamKey + "】的值不能为空")
			}
			yes, err := checkConditions(v, paramValue)
			if err != nil {
				return nil, err
			}
			if yes {
				flag++
			}
		}
		// fmt.Printf("flag=%d\n", flag)
		// 满足所有条件
		if flag == len(n.Properties.Conditions[0]) {
			result = n
		}
	}
	return result, nil
}

func checkConditions(cond *node.ConditionNode, value string) (bool, error) {
	// 判断类型
	switch cond.Type {
	case node.ActionConditionTypes[node.RANGE]:
		val, err := strconv.Atoi(value)
		if err != nil {
			return false, err
		}
		if len(cond.LowerBound) == 0 && len(cond.UpperBound) == 0 && len(cond.LowerBoundEqual) == 0 && len(cond.UpperBoundEqual) == 0 && len(cond.BoundEqual) == 0 {
			return false, errors.New("条件【" + cond.Type + "】的上限或者下限值不能全为空")
		}
		// 判断下限，lowerBound
		if len(cond.LowerBound) > 0 {
			low, err := strconv.Atoi(cond.LowerBound)
			if err != nil {
				return false, err
			}
			if val <= low {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		if len(cond.LowerBoundEqual) > 0 {
			le, err := strconv.Atoi(cond.LowerBoundEqual)
			if err != nil {
				return false, err
			}
			if val < le {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		// 判断上限,upperBound包含等于
		if len(cond.UpperBound) > 0 {
			upper, err := strconv.Atoi(cond.UpperBound)
			if err != nil {
				return false, err
			}
			if val >= upper {
				return false, nil
			}
		}
		if len(cond.UpperBoundEqual) > 0 {
			ge, err := strconv.Atoi(cond.UpperBoundEqual)
			if err != nil {
				return false, err
			}
			if val > ge {
				return false, nil
			}
		}
		if len(cond.BoundEqual) > 0 {
			equal, err := strconv.Atoi(cond.BoundEqual)
			if err != nil {
				return false, err
			}
			if val != equal {
				return false, nil
			}
		}
		return true, nil
	case node.ActionConditionTypes[node.VALUE]:
		if len(cond.ParamValues) == 0 {
			return false, errors.New("条件节点【" + cond.Type + "】的 【paramValues】数组不能为空，值如：'paramValues:['调休','年假']")
		}
		for _, val := range cond.ParamValues {
			if value == val {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("未知的NodeCondition类型【" + cond.Type + "】,正确类型应为以下中的一个:")
	}
}

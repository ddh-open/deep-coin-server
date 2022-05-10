package task

import "devops-http/app/module/base"

// WorkflowTask 流程任务表
// ExecutionID 执行流ID
// Name 任务名称，在流程文件中定义
// TaskDefKey 任务定义的ID值
// Assignee 被指派执行该任务的人
// Owner 任务拥有人
type WorkflowTask struct {
	base.DevopsModel
	// Company 任务创建人对应的公司
	// Company string `json:"company"`
	// ExecutionID     string `json:"executionID"`
	// 当前执行流所在的节点
	NodeID string `json:"nodeId"`
	Step   int    `json:"step"`
	// 流程实例id
	ProcInstID uint   `json:"procInstID"`
	Assignee   string `json:"assignee"`
	ClaimTime  string `json:"claimTime"`
	// 还未审批的用户数，等于0代表会签已经全部审批结束，默认值为1
	MemberCount   int8 `json:"memberCount" gorm:"default:1"`
	UnCompleteNum int8 `json:"unCompleteNum" gorm:"default:1"`
	// 审批通过数
	AgreeNum int8 `json:"agreeNum"`
	// and 为会签，or为或签，默认为or
	ActType    string `json:"actType" gorm:"default:'or'"`
	IsFinished bool   `gorm:"default:false" json:"isFinished"`
}

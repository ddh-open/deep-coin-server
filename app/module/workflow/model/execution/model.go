package execution

import "devops-http/app/module/base"

// WorkflowExecution 流程实例（执行流）表
// ProcInstID 流程实例ID
// BusinessKey 启动业务时指定的业务主键
// WorkflowProcDefID 流程定义数据的ID
type WorkflowExecution struct {
	base.DevopsModel
	Rev         int    `json:"rev"`
	ProcInstID  uint   `json:"procInstID"`
	ProcDefID   uint   `json:"procDefID"`
	ProcDefName string `json:"procDefName"`
	// NodeInfos 执行流经过的所有节点
	NodeInfos string `gorm:"size:4000" json:"nodeInfos"`
	IsActive  int8   `json:"isActive"`
}

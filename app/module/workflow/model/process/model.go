package process

import "devops-http/app/module/base"

// WorkflowInstProc 流程实例
type WorkflowInstProc struct {
	base.DevopsModel
	// 流程定义ID
	ProcDefID uint `json:"procDefId"`
	// 流程定义名
	ProcDefName string `json:"procDefName"`
	// title 标题
	Title string `json:"title"`
	// 用户部门
	Department string `json:"department"`
	Company    string `json:"company"`
	// 当前节点
	NodeID string `json:"nodeID"`
	// 审批人
	Candidate string `json:"candidate"`
	// 当前任务
	TaskID        uint   `json:"taskID"`
	EndTime       string `json:"endTime"`
	Duration      int64  `json:"duration"`
	StartUserID   string `json:"startUserId"`
	StartUserName string `json:"startUserName"`
	IsFinished    bool   `gorm:"default:false" json:"isFinished"`
}

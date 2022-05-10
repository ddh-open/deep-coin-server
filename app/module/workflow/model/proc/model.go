package proc

import "devops-http/app/module/base"

// WorkflowProc 流程定义表
type WorkflowProc struct {
	base.DevopsModel
	Name    string `json:"name,omitempty"`
	Version int    `json:"version,omitempty"`
	// 流程定义json字符串
	Resource string `gorm:"size:10000" json:"resource,omitempty"`
	// 用户id
	UserId   string `json:"userId,omitempty"`
	Username string `json:"username,omitempty"`
	// 用户所在公司
	Company string `json:"company,omitempty"`
}

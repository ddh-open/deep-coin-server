package identity

import "devops-http/app/module/base"

// WorkflowIdentityLink 用户组同任务的关系
type WorkflowIdentityLink struct {
	base.DevopsModel
	Group      string `json:"group,omitempty"`
	Type       string `json:"type,omitempty"`
	UserID     string `json:"userid,omitempty"`
	UserName   string `json:"username,omitempty"`
	TaskID     uint   `json:"taskID,omitempty"`
	Step       int    `json:"step"`
	ProcInstID uint   `json:"procInstID,omitempty"`
	Company    string `json:"company,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// KindIdentity 类型
type KindIdentity int

const (
	// CANDIDATE 候选
	CANDIDATE KindIdentity = iota
	// PARTICIPANT 参与人
	PARTICIPANT
	// MANAGER 上级领导
	MANAGER
	// NOTIFIER 抄送人
	NOTIFIER
)

// KindsIdentity 参与人的类型
var KindsIdentity = [...]string{CANDIDATE: "candidate", PARTICIPANT: "participant", MANAGER: "主管", NOTIFIER: "notifier"}

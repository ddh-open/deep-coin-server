package request

import (
	"devops-http/app/module/workflow/model/node"
)

// ProcRequest 流程定义表
type ProcRequest struct {
	Name string `json:"name"`
	// 流程定义json字符串
	Resource *node.Node `json:"resource"`
	// 用户id
	UserId   string `json:"userId"`
	Username string `json:"username"`
	// 用户所在公司
	Company string `json:"company"`
}

// ReceiverProcess 接收页面传递参数
type ReceiverProcess struct {
	UserID     string             `json:"userId"`
	ProcInstID string             `json:"procInstID"`
	Username   string             `json:"username"`
	Company    string             `json:"company"`
	ProcName   string             `json:"procName"`
	Title      string             `json:"title"`
	Department string             `json:"department"`
	Var        *map[string]string `json:"var"`
}

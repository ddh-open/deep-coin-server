package workflow

import (
	"devops-http/app/module/base/request"
	"devops-http/app/module/workflow/model/proc"
)

// ProcPageReceiver 分页参数
type ProcPageReceiver struct {
	request.PageRequest
	proc.WorkflowProc
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

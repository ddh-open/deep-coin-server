package request

import (
	"devops-http/app/module/workflow/model/proc"
)

// ProcPageReceiver 分页参数
type ProcPageReceiver struct {
	PageRequest
	proc.WorkflowProc
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

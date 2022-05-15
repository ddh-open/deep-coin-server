package request

import (
	"devops-http/app/module/workflow/model/process"
)

// ProcessPageReceiver 分页参数
type ProcessPageReceiver struct {
	PageRequest
	process.WorkflowInstProc
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

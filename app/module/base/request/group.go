package request

import (
	"devops-http/app/module/base"
)

// SearchGroupParams group的筛选
type SearchGroupParams struct {
	base.DevopsSysGroup
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

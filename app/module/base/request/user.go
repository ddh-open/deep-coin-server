package request

import (
	"devops-http/app/module/base"
)

// SearchUserParams user的筛选
type SearchUserParams struct {
	base.DevopsSysUser
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

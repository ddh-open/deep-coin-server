package request

import "devops-http/app/module/cluster/model/application"

// SearchApplicationParams 查询服务
type SearchApplicationParams struct {
	application.DevopsClusterApplication
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

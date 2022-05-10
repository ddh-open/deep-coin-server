package request

import (
	"devops-http/app/module/sys/model/config"
	"devops-http/app/module/sys/model/menu"
	"devops-http/app/module/sys/model/operation"
	"devops-http/app/module/sys/model/path"
	"devops-http/app/module/sys/model/role"
)

// SearchApiParams api分页条件查询及排序结构体
type SearchApiParams struct {
	path.DevopsSysApi
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// SearchConfigParams config的筛选
type SearchConfigParams struct {
	config.DevopsSysConfig
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// SearchMenusParams menus的筛选
type SearchMenusParams struct {
	menu.DevopsSysMenu
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// SearchLogsParams 操作记录的筛选
type SearchLogsParams struct {
	operation.DevopsSysOperationRecord
	PageRequest
	OrderKey   string   `json:"orderKey"` // 排序
	Desc       bool     `json:"desc"`     // 排序方式:升序false(默认)|降序true
	TimeFilter []string `json:"timeFilter"`
}

// CopyRoleParams role的复制
type CopyRoleParams struct {
	role.DevopsSysRole
	CopyId string `json:"copyId"`
}

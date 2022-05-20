package request

import (
	"devops-http/app/module/sys/model/config"
	"devops-http/app/module/sys/model/host"
	"devops-http/app/module/sys/model/icon"
	"devops-http/app/module/sys/model/menu"
	"devops-http/app/module/sys/model/operation"
	"devops-http/app/module/sys/model/path"
	"devops-http/app/module/sys/model/role"
)

// SearchHostParams 查询主机
type SearchHostParams struct {
	host.DevopsCmdbHostGroup
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// SearchIconParams 图标分页条件查询及排序结构体
type SearchIconParams struct {
	icon.DevopsSysIcon
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// SearchRoleParams 角色分页条件查询及排序结构体
type SearchRoleParams struct {
	role.DevopsSysRole
	PageRequest
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

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

type LoginRequest struct {
	Username string
	Password string
	Type     int
}

type ChangePasswordRequest struct {
	Username    string
	Password    string
	OldPassword string
	Type        int
}

type UserRelativeRoleRequest struct {
	UserId  string   `json:"userId"`
	RoleIds []string `json:"roleIds"`
}

type GroupRelativeUserRequest struct {
	GroupId int   `json:"groupId"`
	UserIds []int `json:"userIds"`
}

type RelativeRoleMenuRequest struct {
	RoleId  string   `json:"roleId"`
	MenuIds []string `json:"menuIds"`
}

type RelativeRoleApisRequest struct {
	RoleId string   `json:"roleId"`
	ApiIds []string `json:"apiIds"`
}

// GetByRoleId Get role by id structure
type GetByRoleId struct {
	RoleId string `json:"roleId"` // 角色ID
}

type DeleteById struct {
	Id int `json:"ids"` // 角色ID
}

type ReqById struct {
	Ids string `json:"ids"` // 角色ID
}

type Empty struct{}

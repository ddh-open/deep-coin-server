package user

import "devops-http/app/module/base"

type DevopsSysUserEntity struct {
	base.DevopsSysUser
	RoleIds []string `json:"roleIds"`
}

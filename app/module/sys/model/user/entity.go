package user

type DevopsSysUserEntity struct {
	DevopsSysUser
	RoleIds []string `json:"roleIds"`
}

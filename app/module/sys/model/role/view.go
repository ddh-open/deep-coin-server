package role

// DevopsSysRoleView 角色表
type DevopsSysRoleView struct {
	DevopsSysRole
	DomainName string `json:"domainName"`
}

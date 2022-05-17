package role

// DevopsSysRoleView 角色表
type DevopsSysRoleView struct {
	DevopsSysRole
	Menus      []int  `json:"menus"`
	Apis       []int  `json:"apis"`
	DomainName string `json:"domainName"`
}

package sys

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

type RelativeUserRequest struct {
	UserId  string   `json:"userId"`
	RoleIds []string `json:"roleIds"`
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
	Ids string `json:"ids"` // 角色ID
}

type RequestById struct {
	Ids string `json:"ids"` // 角色ID
}

type Empty struct{}

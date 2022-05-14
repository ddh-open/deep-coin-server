package user

import "devops-http/app/module/base"

// DevopsSysUserView 用户
type DevopsSysUserView struct {
	DevopsSysUser
	RoleIds []string `json:"roleIds"`
	Roles   []string `json:"roles"`
}

// SysUserViewColumns get sql column name.获取数据库列名
var SysUserViewColumns = struct {
	ID        base.DevopsColumns `json:"id"`
	UUID      base.DevopsColumns `json:"uuid"`
	Username  base.DevopsColumns `json:"username"`
	Password  base.DevopsColumns `json:"-"`
	Salt      base.DevopsColumns `json:"-"`
	RealName  base.DevopsColumns `json:"realName"`
	UserType  base.DevopsColumns `json:"userType"`
	Status    base.DevopsColumns `json:"status"`
	WorkNum   base.DevopsColumns `json:"workNum"`
	EndTime   base.DevopsColumns `json:"-"`
	Email     base.DevopsColumns `json:"email"`
	Tel       base.DevopsColumns `json:"tel"`
	Address   base.DevopsColumns `json:"address"`
	TitleURL  base.DevopsColumns `json:"-"`
	Remark    base.DevopsColumns `json:"remark"`
	Theme     base.DevopsColumns `json:"theme"`
	Enable    base.DevopsColumns `json:"enable"`
	UpdatedAt base.DevopsColumns `json:"updatedAt"`
	UpdateID  base.DevopsColumns `json:"updateId"`
	CreatedAt base.DevopsColumns `json:"createdAt"`
	CreateID  base.DevopsColumns `json:"-"`
	DeletedAt base.DevopsColumns `json:"-"`
	RoleIds   base.DevopsColumns `json:"roleIds"`
}{
	ID:        base.DevopsColumns{Width: "60", Label: "id", Prop: "id"},
	UUID:      base.DevopsColumns{Width: "200", Label: "uuid", Prop: "uuid"},
	Username:  base.DevopsColumns{Width: "100", Label: "用户名", Prop: "username"},
	RealName:  base.DevopsColumns{Width: "100", Label: "真实名", Prop: "realName"},
	UserType:  base.DevopsColumns{Width: "100", Label: "用户类型", Prop: "userType"},
	Status:    base.DevopsColumns{Width: "100", Label: "状态", Prop: "status"},
	WorkNum:   base.DevopsColumns{Width: "100", Label: "工号", Prop: "workNum"},
	EndTime:   base.DevopsColumns{Width: "100", Label: "上次登录时间", Prop: "endTime"},
	Email:     base.DevopsColumns{Width: "100", Label: "邮件", Prop: "email"},
	Tel:       base.DevopsColumns{Width: "100", Label: "电话", Prop: "tel"},
	Address:   base.DevopsColumns{Width: "100", Label: "地址", Prop: "address"},
	TitleURL:  base.DevopsColumns{Width: "100", Label: "头像", Prop: "titleURL"},
	Remark:    base.DevopsColumns{Width: "100", Label: "描述", Prop: "remark"},
	Theme:     base.DevopsColumns{Width: "100", Label: "主题", Prop: "theme"},
	Enable:    base.DevopsColumns{Width: "100", Label: "是否启动", Prop: "enable"},
	UpdatedAt: base.DevopsColumns{Width: "100", Label: "更新时间", Prop: "updatedAt"},
	UpdateID:  base.DevopsColumns{Width: "100", Label: "更新者", Prop: "updateID"},
	CreatedAt: base.DevopsColumns{Width: "100", Label: "创建时间", Prop: "createdAt"},
	RoleIds:   base.DevopsColumns{Width: "100", Label: "角色", Prop: "roleIds"},
}

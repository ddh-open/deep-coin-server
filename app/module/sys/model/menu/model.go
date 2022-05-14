package menu

import "devops-http/app/module/base"

// DevopsSysMenu 菜单表
type DevopsSysMenu struct {
	base.DevopsModel
	ParentId  string                            `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path      string                            `json:"path" gorm:"comment:路由path"`        // 路由path
	Name      string                            `json:"name" gorm:"comment:路由name"`        // 路由name
	Component string                            `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Redirect  string                            `json:"redirect" gorm:"comment:重定向"`       // 排序标记
	Sort      int                               `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	Title     string                            `json:"title" gorm:"comment:菜单名"`          // 菜单名
	Meta      `json:"meta" gorm:"comment:附加属性"` // 附加属性
	Children  []DevopsSysMenu                   `json:"children" gorm:"-"`
}

type Meta struct {
	NoKeepAlive bool   `json:"noKeepAlive" gorm:"comment:是否缓存"`         // 是否缓存
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"` // 是否是基础路由（开发中）
	Icon        string `json:"icon" gorm:"comment:菜单图标"`                // 菜单图标
	NoClosable  bool   `json:"noClosable" gorm:"comment:是否固定"`          // 是否固定
	LevelHidden bool   `json:"levelHidden" gorm:"comment:始终显示当前节点"`     // 始终显示当前节点
	Hidden      bool   `json:"hidden" gorm:"comment:是否隐藏"`              // 是否隐藏
	Dot         bool   `json:"dot" gorm:"comment:是否dot"`                // 是否dot
	Badge       string `json:"badge" gorm:"comment:badge"`              // badge
}

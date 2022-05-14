package icon

import "devops-http/app/module/base"

// DevopsSysIcon 系统图标
type DevopsSysIcon struct {
	base.DevopsModel
	Title string `json:"title" gorm:"comment:系统图标名称"` // 系统图标名称
}

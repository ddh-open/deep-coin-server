package user

import "devops-http/app/module/base"

// DevopsSysCabinRule 鉴权规则
type DevopsSysCabinRule struct {
	base.DevopsModel
	Source   string `json:"source" gorm:"column:source"`
	Resource string `json:"resource" gorm:"column:resource"`
	Domain   string `json:"domain" gorm:"column:domain"`
	Method   string `json:"method" gorm:"column:method"`
}

// CabinOutInfo Cabin info structure
type CabinOutInfo struct {
	PType    string `json:"pType"`
	Resource string `json:"resource"` // 资源
	Domain   string `json:"domain"`   // 方法
	Method   string `json:"method"`   // 方法
}

// CabinInReceive Cabin structure for input parameters
type CabinInReceive struct {
	PType    string `json:"pType"`    //
	Source   string `json:"source"`   //
	Resource string `json:"resource"` // 路径
	Domain   string `json:"domain"`   // 方法
	Method   string `json:"method"`   // 方法
}
